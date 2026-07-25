package sqlite

import (
	"context"
	"database/sql"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/transports"
	"fmt"
	"log"
	"net/url"
	_"modernc.org/sqlite"
)

type Connecion struct {
	config entities.ConnectionConfig

	transport transports.Transport

	db        *sql.DB
	connected bool
}

func New(config entities.ConnectionConfig, transport transports.Transport) *Connecion {
	return &Connecion{
		config:    config,
		transport: transport,
		connected: false,
	}
}

func (c *Connecion) ID() string {
	return c.config.ID
}

func (c *Connecion) Name() string {
	return c.config.Name
}

func (c *Connecion) DatabaseName() string {
	return c.config.Database
}

func (c *Connecion) Type() string {
	return "sqlite"
}

func (c *Connecion) DB() *sql.DB {
	return c.db
}

func (c *Connecion) dsn() string {
	u, err := url.Parse("file:" + c.config.Database)
	if err != nil {
		return "file:" + c.config.Database + "?_journal_mode=WAL&_foreign_keys=on"
	}
	q := u.Query()
	q.Set("_journal_mode", "WAL")
	q.Set("_foreign_keys", "on")
	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Connecion) Connect(ctx context.Context) error {
	if err := c.transport.Connect(ctx); err != nil {
		return err
	}
	log.Println("connected to transport")
	log.Println("opening database connection")
	log.Println("dsn:", c.dsn())
	db, err := sql.Open("sqlite", c.dsn())
	if err != nil {
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return err
	}

	c.db = db
	c.connected = true

	return nil
}

func (c *Connecion) Disconnect() error {
	if !c.connected {
		return nil
	}

	if c.db == nil {
		return nil
	}

	err := c.db.Close()
	c.connected = false
	return err
}

func (c *Connecion) Ping(ctx context.Context) error {
	if !c.connected {
		return fmt.Errorf("not connected")
	}
	if c.db == nil {
		return fmt.Errorf("db not initialized")
	}
	return c.db.PingContext(ctx)
}

func (c *Connecion) IsConnected() bool {
	return c.connected
}
