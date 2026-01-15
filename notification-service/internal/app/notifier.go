package app

import "context"

type Notifier interface {
	Notify(ctx context.Context, n Notification) error
}

type Notification struct {
	Text string
}
