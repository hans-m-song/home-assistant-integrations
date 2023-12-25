package broker

import (
	"context"
	"time"
)

type Broker interface {
	Name() string
	Publish(topic string, payload any) error
	Subscribe(topic string, handler func(topic string, payload any)) error
	Health(ctx context.Context) (map[string]any, error)
	Disconnect(deadline time.Duration) error
}
