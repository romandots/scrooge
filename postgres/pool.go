package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"scrooge/config"
)

func InitPool() (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Database,
	)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	return pgxpool.NewWithConfig(ctx, config)
}
