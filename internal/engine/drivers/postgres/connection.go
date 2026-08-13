package postgres

import (
	"context"
	"database/sql"
	"db-viewer/internal/engine/entities"
	"db-viewer/internal/engine/transports"
	"fmt"
	"log"
	"strings"

	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	return string(entities.DialectPostgreSQL)
}

func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) dsn() string {
	// Fallback to config details if transport provides a non-standard network address
	host := c.transport.Address()
	if host == "localfs" || host == "" {
		host = c.config.Host
	}
	parts := []string{
		fmt.Sprintf("host=%s", host),
		fmt.Sprintf("port=%d", c.config.Port),
		fmt.Sprintf("user=%s", c.config.User),
		fmt.Sprintf("dbname=%s", c.config.Database),
		"sslmode=disable",
	}

	if c.config.Password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", c.config.Password))
	}

	return strings.Join(parts, " ")
}

func (c *Connection) Connect(ctx context.Context) error {
	if err := c.transport.Connect(ctx); err != nil {
		return err
	}
	log.Println("conneccting to postgres pgx with dsn ", c.dsn())

	db, err := sql.Open("pgx", c.dsn())
	if err != nil {
		return err
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return err
	}

	c.db = db
	c.connected = true

	var database, user, server, version string

	err = db.QueryRow(`
    SELECT
        current_database(),
        current_user,
        inet_server_addr()::text,
        version()
	`).Scan(&database, &user, &server, &version)

	if err != nil {
		return fmt.Errorf("connection verification failed: %w", err)
	}

	fmt.Printf("DATABASE: %s\n", database)
	fmt.Printf("USER: %s\n", user)
	fmt.Printf("SERVER: %s\n", server)
	fmt.Printf("VERSION: %s\n", version)

	return nil
}

func (c *Connection) Disconnect() error {
	if c.db == nil {
		return nil
	}

	err := c.db.Close()
	c.connected = false
	return err
}

func (c *Connection) Ping(ctx context.Context) error {
	if c.db == nil {
		return fmt.Errorf("postgres connection not initialized")
	}

	return c.db.PingContext(ctx)
}

func (c *Connection) IsConnected() bool {
	return c.connected
}

func (c *Connection) Config() entities.ConnectionConfig {
	return c.config
}