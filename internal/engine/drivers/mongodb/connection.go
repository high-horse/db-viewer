package mongodb

import (
	"context"
	"crypto/tls"
	"db-lens/internal/engine/entities"
	"db-lens/internal/engine/transports"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Connection struct {
	config    entities.ConnectionConfig
	transport transports.Transport
	client    *mongo.Client
}

func New(config entities.ConnectionConfig, transport transports.Transport) *Connection {
	return &Connection{config: config, transport: transport}
}
func (c *Connection) ID() string                        { return c.config.ID }
func (c *Connection) Name() string                      { return c.config.Name }
func (c *Connection) Type() string                      { return "mongodb" }
func (c *Connection) DatabaseName() string              { return c.config.Database }
func (c *Connection) Config() entities.ConnectionConfig { return c.config }
func (c *Connection) Client() *mongo.Client             { return c.client }
func (c *Connection) DB() *mongo.Database {
	if c.client == nil {
		return nil
	}
	return c.client.Database(c.config.Database)
}
func (c *Connection) IsConnected() bool { return c.client != nil }

func isURI(host string) bool {
	return strings.HasPrefix(host, "mongodb://") || strings.HasPrefix(host, "mongodb+srv://")
}

func (c *Connection) clientOptions() (*options.ClientOptions, error) {
	opts := options.Client()
	if isURI(c.config.Host) {
		if c.config.SSHConfig != nil {
			return nil, fmt.Errorf("use a database host and port instead of a MongoDB URI with SSH")
		}
		opts.ApplyURI(c.config.Host)
	} else {
		opts.SetHosts([]string{c.transport.Address()})
	}
	if c.config.User != "" {
		opts.SetAuth(options.Credential{Username: c.config.User, Password: c.config.Password, AuthSource: "admin"})
	}
	if c.config.SSL {
		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
		if c.config.SSHConfig != nil {
			tlsConfig.ServerName = c.config.Host
		}
		opts.SetTLSConfig(tlsConfig)
	}
	if c.config.SSHConfig != nil {
		opts.SetDirect(true)
	}
	opts.SetConnectTimeout(15 * time.Second).SetServerSelectionTimeout(15 * time.Second)
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid MongoDB connection: %w", err)
	}
	return opts, nil
}

func (c *Connection) Connect(ctx context.Context) error {
	if c.client != nil {
		return nil
	}
	if strings.TrimSpace(c.config.Database) == "" {
		return fmt.Errorf("MongoDB database is required")
	}
	// A URI is handled by the MongoDB driver; Direct.Connect has no network side effects.
	if err := c.transport.Connect(ctx); err != nil {
		return err
	}
	opts, err := c.clientOptions()
	if err != nil {
		_ = c.transport.Close()
		return err
	}
	setup, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(setup, opts)
	if err == nil {
		err = client.Ping(setup, readpref.Primary())
	}
	if err != nil {
		if client != nil {
			cleanup, done := context.WithTimeout(context.Background(), 5*time.Second)
			_ = client.Disconnect(cleanup)
			done()
		}
		_ = c.transport.Close()
		return fmt.Errorf("connect to MongoDB: %w", err)
	}
	c.client = client
	return nil
}
func (c *Connection) Ping(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("MongoDB is not connected")
	}
	ping, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	return c.client.Ping(ping, readpref.Primary())
}
func (c *Connection) Disconnect() error {
	var err error
	if c.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = c.client.Disconnect(ctx)
		cancel()
		c.client = nil
	}
	return errors.Join(err, c.transport.Close())
}
