package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/simpler-tha/internal/config"
)

type Client struct {
	Conn *pgx.Conn
}

// NewClient sets up the Postgres client.
func NewClient(ctx context.Context, cfg config.Postgres) (*Client, error) {
	connStr := buildConnString(cfg.Username, cfg.Pass, cfg.Host, cfg.Database, cfg.Port)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Client{
		Conn: conn,
	}, nil
}

func buildConnString(user, pass, host, dbName string, port int) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, dbName)
}
