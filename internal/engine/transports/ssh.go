package transports

import (
	"context"
	"fmt"
)

type SSH struct {
	dbHost string
	dbPort int

	sshHost     string
	sshPort     int
	sshUser     string
	sshKeyPath  string
	sshPassword string
}

func NewSSH(dbHost string, dbPort int, sshHost string, sshPort int, sshUser string, sshKeyPath string, sshPassword string) *SSH {
	return &SSH{
		dbHost:      dbHost,
		dbPort:      dbPort,
		sshHost:     sshHost,
		sshPort:     sshPort,
		sshUser:     sshUser,
		sshKeyPath:  sshKeyPath,
		sshPassword: sshPassword,
	}
}

func (d *SSH) Connect(ctx context.Context) error {
	// TODO: Future implementation:
	// 1. Dial sshHost:sshPort using x/crypto/ssh
	// 2. Establish a background tunnel listener targeting dbHost:dbPort

	return nil
}

func (d *SSH) Close() error {
	return nil
}

func (d *SSH) Address() string {
	return fmt.Sprintf("%s:%d", d.sshHost, d.sshPort)
}
