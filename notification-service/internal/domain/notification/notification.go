package notification

import "context"

type Notification struct {
	Text string
}

type Notifier interface {
	Notify(ctx context.Context, n Notification) error
}
