package notifications

type Notifier interface {
	SendMessage(msg string) error
}
