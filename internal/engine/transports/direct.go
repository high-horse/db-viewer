package transports

import (
	"context"
	"net"
	"strconv"
)

type Direct struct {
	host string
	port int
}

func NewDirect(host string, port int) *Direct {
	return &Direct{
		host: host,
		port: port,
	}
}

func (d *Direct) Connect(ctx context.Context) error {
	return nil
}

func (d *Direct) Close() error {
	return nil
}

func (d *Direct) Address() string {
	host := d.host
	if host == "" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, strconv.Itoa(d.port))
}
