package pulse

// defineStatus is a helper function to define status message.
func defineStatus(status string) string {
	switch status {
	case StatusConfirmed:
		return MessageConfirmed
	case StatusPending:
		return MessagePending
	case StatusFailed:
		return MessageFailed
	case StatusDNE:
		return MessageDNE
	default:
		return "Unknown"
	}
}
