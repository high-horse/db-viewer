package transports

import (
	"context"
	"fmt"
)

type Direct struct{
	host string
	port int
}

func NewDirect(host string, port int) *Direct {
	return &Direct{
		host: host,
		port: port,
	}
}

func(d *Direct) Connect(ctx context.Context) error {
	return nil
}

func(d *Direct) Close() error {
	return nil
}

func(d *Direct) Address() string {
	return fmt.Sprintf("%s:%d", d.host, d.port)
}