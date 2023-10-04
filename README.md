# go-pulse

Broadcast, monitor, blockchain transaction in Go.

## Installation

### Prerequisites

Make sure you have **[Go](https://go.dev/)** installed, preferably one of the **three latest major versions**.

### Installation

You can install Pulse using `go get`:

```
go get github.com/markpassawat/go-pulse
```

Alternatively, If you're using Go module support, you can add the following import statement to your code:

```
import "github.com/markpassawat/go-pulse"
```

## Quick start

You can explore a variety of ready-to-run examples in the [Pulse example repository](https://github.com/markpassawat/go-pulse/blob/main/example/basic/main.go).

### Broadcasting and Monitoring a Transaction

Here's an example in Go demonstrating how to broadcast a transaction and subsequently monitor it using Pulse:

```go
func main() {
	// Create a Pulse client with the default HTTP server for interaction with.
	p := pulse.NewPulse()

	// Use the BroadcastAndMonitor function to broadcast and monitor the transaction.
	tx, err := p.BroadcastAndMonitor(&pulse.Asset{
		Symbol:    "ETH",            // Transaction symbol (string)
		Price:     4500,             // Symbol price (uint64)
		Timestamp: 1678912345,       // Unix timestamp (uint64)
	})
	if err != nil {
		// Handle the error here.
	}

	// Output:
	// tx.TxHash: 131cc9cb7ee149a4a629d9fed28290026caf21bbb54cee8e661375ae9b124376
	// tx.Status: CONFIRMED
	// tx.Message: Transaction has been processed and confirmed
}

```

### Broadcasting and Monitoring Multiple Transactions

Here's an example in Go demonstrating how to broadcast and subsequently monitor multiple transactions using Pulse:

```go
func main() {
	// Create a Pulse client with the default HTTP server for interaction.
	p := pulse.NewPulse()

	// Define a list of assets to be broadcasted and monitored.
	assets := []pulse.Asset{
		{
			Symbol:    "ETH",
			Price:     4500,
			Timestamp: 1678912345,
		},
		{
			Symbol:    "BTC",
			Price:     4500,
			Timestamp: 1678912345,
		},
	}

	// Use the MultipleBroadcastAndMonitor function
	// to broadcast and monitor the assets.
	txs, err := p.MultipleBroadcastAndMonitor(assets)
	if err != nil {
		// Handle the error here.
	}
}

```

## Strategies for Managing Transaction Status

Explore the following strategies for controlling different transaction statuses, and consider persistently store information logs in a database for future use:

`CONFIRMED`, the transaction has been processed and confirmed. Integrate with notify service to notify that broadcasting transaction was success. Notify that the transaction was successful.

`FAILED`, the transaction failed to process. Integrate with notify service to notify that broadcasted transaction was not success. It's possible that the transaction was rejected by the server or the transaction was rejected by the network (e.g., insufficient funds, invalid transaction, etc.). Prompt it to the user to try again.

`DNE (Does Not Exist)`, Transaction does not exist. Integrate with notify service to notify that the broadcasted transaction was does not exist which mean it wasn't processed or recorded by the server. It's possible that provided an incorrect hash or that the transaction was never broadcasted to the network.

By following these strategies and storing transaction-related data in a database, you can maintain a clear record of transaction statuses, troubleshoot issues effectively, and keep stakeholders informed about the outcome of their transactions.
