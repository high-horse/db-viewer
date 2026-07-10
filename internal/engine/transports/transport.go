package transports

import "context"

type Transport interface {
	Connect(ctx context.Context) error
	Close() error

	Address() string
}