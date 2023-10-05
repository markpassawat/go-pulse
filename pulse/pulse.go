package pulse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Asset struct {
	// Symbol is the symbol of the asset.
	Symbol string `json:"symbol"`
	// Price is the price of the asset.
	Price uint64 `json:"price"`
	// Timestamp is the timestamp price retrieval.
	Timestamp uint64 `json:"timestamp"`
}

// Validate validates the asset data request.
func (a Asset) Validate() error {
	return v.ValidateStruct(&a,
		v.Field(&a.Symbol, v.Required),
		v.Field(&a.Price, v.Required),
		v.Field(&a.Timestamp, v.Required),
	)
}

type broadcastResponse struct {
	// TxHash is the transaction hash returned from broadcast the asset.
	TxHash string `json:"tx_hash"`
}

type txHashData struct {
	// TxHash is the transaction hash.
	TxHash string `json:"tx_hash"`
	// Status is the status of the transaction (CONFIRMED, PENDING, FAILED, DNE).
	Status string `json:"tx_status"`
	// Message is the message to describe the status.
	Message string `json:"message"`
}

type txStatusResponse struct {
	// Status is the status of the transaction (CONFIRMED, PENDING, FAILED, DNE).
	Status string `json:"tx_status"`
}

type Options struct {
	// URL is the URL for broadcast asset data and monitor its status.
	URL string
	// PulseInterval is the time interval for monitor status.
	PulseInterval time.Duration
}

// DefaultURL is the default URL for broadcast asset data and monitor its status.
const DefaultURL = "https://mock-node-wgqbnxruha-as.a.run.app"

// DefaultPulseInterval is the default time interval for monitor status.
const DefaultPulseInterval = 3 * time.Second

// Status and definition for the transaction status.
var (
	StatusConfirmed  = "CONFIRMED"
	MessageConfirmed = "Transaction has been processed and confirmed"
	StatusPending    = "PENDING"
	MessagePending   = " Transaction is awaiting processing"
	StatusFailed     = "FAILED"
	MessageFailed    = "Transaction failed to process"
	StatusDNE        = "DNE"
	MessageDNE       = "Transaction does not exist"
)

// pulse is the client for pulse,
// it should be created via New or NewPulse.
// NewPulse will use the default URL and pulse interval.
// New will use the URL and pulse interval from the options.
type pulse struct {
	client *resty.Client
	pulse  time.Duration
}

// Ensure pulse implements IPulse.
var _ IPulse = new(pulse)

// New returns a new pulse client with options
func New(opts *Options) *pulse {
	if opts.URL == "" {
		return NewPulse()
	}
	if opts.PulseInterval == 0 {
		opts.PulseInterval = DefaultPulseInterval
	}
	client := resty.New().SetBaseURL(opts.URL)
	return &pulse{
		client: client,
		pulse:  opts.PulseInterval,
	}
}

// NewPulse returns a new pulse client with default options.
func NewPulse() *pulse {
	return &pulse{
		client: resty.New().SetBaseURL(DefaultURL),
		pulse:  DefaultPulseInterval,
	}
}

// Broadcast broadcasts the asset data to the server.
// It returns the transaction hash and error if any.
func (p *pulse) Broadcast(asset *Asset) (*broadcastResponse, error) {
	if err := asset.Validate(); err != nil {
		return nil, fmt.Errorf("validate failed broadcast_asset: %w", err)
	}

	var result broadcastResponse
	resp, err := p.client.R().
		SetBody(Asset{
			Symbol:    asset.Symbol,
			Price:     asset.Price,
			Timestamp: asset.Timestamp,
		}).SetResult(&result).Post("/broadcast")
	if err != nil {
		return nil, fmt.Errorf("request failed broadcast_asset: %w", err)
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return nil, fmt.Errorf("broadcast failed: %v", err)
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("unmarshal failed broadcast_asset: %w", err)
	}

	return &result, nil
}

// MultipleBroadcast broadcasts multiple asset data to the server.
// It returns  multiple transaction hash and error if any of each request.
func (p *pulse) MultipleBroadcast(opts []Asset) ([]broadcastResponse, error) {
	var result []broadcastResponse
	for _, v := range opts {
		hash, err := p.Broadcast(&v)
		if err != nil {
			return nil, fmt.Errorf("request failed multiple_broadcast: %w", err)
		}
		result = append(result, *hash)
	}
	return result, nil
}

// MonitorStatus monitors the status of the transaction hash that obtained from Broadcast.
// It returns the transaction hash data and error if any.
// It will keep monitoring the status until the status is not pending.
// It will sleep for the pulse interval before each request.
func (p *pulse) MonitorStatus(txHash string) (*txHashData, error) {
	for {
		var result txStatusResponse
		resp, err := p.client.R().SetPathParams(map[string]string{
			"tx_hash": txHash,
		}).SetResult(&result).Get("/check/{tx_hash}")
		if err != nil {
			return nil, fmt.Errorf("request failed monitor_status: %w", err)
		}
		if resp.StatusCode() >= http.StatusBadRequest {
			return nil, fmt.Errorf("monitor status failed: %v", err)
		}

		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			return nil, fmt.Errorf("unmarshal failed monitor_status: %w", err)
		}

		if result.Status != StatusPending {
			return &txHashData{
				TxHash:  txHash,
				Status:  result.Status,
				Message: defineStatus(result.Status),
			}, nil
		}

		// Sleep for the pulse interval before each request
		// to avoid too many requests to the server.
		// The default pulse interval is 3 seconds.
		time.Sleep(p.pulse)
	}
}

// MultipleMonitorStatus monitors the status of multiple transaction hash that obtained from Broadcast.
// It returns the multiple transaction hash data.
func (p *pulse) MultipleMonitorStatus(txsHash ...string) []txHashData {
	txsHashData := []txHashData{}
	for _, hash := range txsHash {
		status, err := p.MonitorStatus(hash)
		if err != nil {
			logrus.Warnf("monitor status of tx hash %v failed: %v\n", txsHash, err)
			continue // Skip this iteration and move to the next hash
		}
		txsHashData = append(txsHashData, txHashData{
			TxHash:  hash,
			Status:  status.Status,
			Message: status.Message,
		})
	}

	return txsHashData
}

// BroadcastAndMonitor broadcasts the asset data to the server and monitor its status.
func (p *pulse) BroadcastAndMonitor(asset *Asset) (*txHashData, error) {
	if err := asset.Validate(); err != nil {
		return nil, fmt.Errorf("validate failed broadcast_and_monitor: %w", err)
	}

	var result txHashData
	resp, err := p.client.R().
		SetBody(Asset{
			Symbol:    asset.Symbol,
			Price:     asset.Price,
			Timestamp: asset.Timestamp,
		}).SetResult(&result).Post("/broadcast")
	if err != nil {
		return nil, fmt.Errorf("request failed broadcast_and_monitor: %w", err)
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return nil, fmt.Errorf("broadcast failed: %v", err)
	}

	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("unmarshal failed broadcast_and_monitor: %w", err)
	}

	status, err := p.MonitorStatus(result.TxHash)
	if err != nil {
		return nil, fmt.Errorf("request failed broadcast_and_monitor: %w", err)
	}
	result.Status = status.Status
	result.Message = status.Message

	return &result, nil
}

// MultipleBroadcastAndMonitor broadcasts multiple asset data to the server and monitor its status.
func (p *pulse) MultipleBroadcastAndMonitor(assets []Asset) ([]txHashData, error) {
	var result []txHashData
	for _, v := range assets {
		hash, err := p.BroadcastAndMonitor(&v)
		if err != nil {
			return nil, fmt.Errorf("request failed multiple_broadcast_and_monitor: %w", err)
		}
		result = append(result, *hash)
	}
	return result, nil
}
