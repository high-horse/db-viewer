package mysql

import (
	"context"
	"database/sql"
	"db-lens/internal/engine/entities"
	"db-lens/internal/engine/transports"
	"errors"
	"fmt"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
)

type Connection struct {
	config entities.ConnectionConfig

	transport transports.Transport

	db        *sql.DB
	connected bool
}

func New(config entities.ConnectionConfig, transport transports.Transport) *Connection {
	return &Connection{
		config:    config,
		transport: transport,
	}
}

func (c *Connection) ID() string {
	return c.config.ID
}

func (c *Connection) Name() string {
	return c.config.Name
}

func (c *Connection) DatabaseName() string {
	return c.config.Database
}

func (c *Connection) Type() string {
	return string(entities.DialectMySQL)
}

func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) dsn() string {
	cfg := mysqldriver.NewConfig()
	cfg.User, cfg.Passwd, cfg.DBName = c.config.User, c.config.Password, c.config.Database
	cfg.Net, cfg.Addr = "tcp", c.transport.Address()
	cfg.ParseTime = true
	cfg.Timeout = 15 * time.Second
	if c.config.SSL {
		cfg.TLSConfig = "true"
	}
	return cfg.FormatDSN()
}

func (c *Connection) Connect(ctx context.Context) error {
	if err := c.transport.Connect(ctx); err != nil {
		return err
	}

	db, err := sql.Open(
		"mysql",
		c.dsn(),
	)

	if err != nil {
		_ = c.transport.Close()
		return err
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		_ = c.transport.Close()
		return err
	}

	c.db = db
	c.connected = true

	return nil
}

func (c *Connection) Disconnect() error {
	var err error
	if c.db != nil {
		err = c.db.Close()
		c.db = nil
	}
	c.connected = false
	return errors.Join(err, c.transport.Close())
}

func (c *Connection) Ping(ctx context.Context) error {

	if !c.connected {
		return fmt.Errorf("not connected")
	}

	if c.db == nil {
		return fmt.Errorf("mysql connection not initialized")
	}

	return c.db.PingContext(ctx)
}

func (c *Connection) IsConnected() bool {
	return c.connected
}

func (c *Connection) Config() entities.ConnectionConfig {
	return c.config
}
