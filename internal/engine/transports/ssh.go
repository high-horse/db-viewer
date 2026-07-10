package transports

import (
	"context"
	"fmt"
)

type SSH struct{
	host string
	port int
}

func NewSSH(host string, port int) *SSH {
	return &SSH{
		host: host,
		port: port,
	}
}

func(d *SSH) Connect(ctx context.Context) error {
	return nil
}

func(d *SSH) Close() error {
	return nil
}

func(d *SSH) Address() string {
	return fmt.Sprintf("%s:%d", d.host, d.port)
}