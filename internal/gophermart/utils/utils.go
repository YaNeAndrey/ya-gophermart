package utils

import (
	"context"
	"database/sql"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TryToOpenDBConnection(dbConnectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbConnectionString)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//var bufError error
	err = retry.Retry(

		func(attempt uint) error {
			if err = db.PingContext(ctx); err != nil {
				return err
			}

			return nil
		},
		strategy.Limit(4),
		strategy.Backoff(backoff.Incremental(-1*time.Second, 2*time.Second)),
	)

	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
