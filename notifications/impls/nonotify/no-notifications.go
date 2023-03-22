package nonotify

type NoNotification struct {
}

func (nn *NoNotification) SendMessage(msg string) error {
	return nil
}
