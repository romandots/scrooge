package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"scrooge/config"
)

var Pool *pgxpool.Pool
var Sq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func InitPool() error {
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
		return err
	}
	ctx := context.Background()

	Pool, err = pgxpool.NewWithConfig(ctx, config)

	return err
}

func exec(sql string, args ...interface{}) error {
	_, err := Pool.Exec(context.Background(), sql, args...)
	return err
}

func query(sql string, args ...interface{}) (pgx.Rows, error) {
	rows, err := Pool.Query(context.Background(), sql, args...)
	return rows, err
}
