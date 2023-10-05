package pulse

type IPulse interface {
	// Broadcast broadcasts the asset data to the server.
	Broadcast(asset *Asset) (*broadcastResponse, error)
	// MonitorStatus monitors the status of the transaction hash that obtained from broadcasting.
	MonitorStatus(txHash string) (*txHashData, error)
	// BroadcastAndMonitor broadcasts the asset data to the server keeps track of its completed state.
	BroadcastAndMonitor(asset *Asset) (*txHashData, error)
}
