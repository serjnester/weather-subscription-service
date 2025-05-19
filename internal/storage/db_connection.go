package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/serjnester/weather-subscription-service/internal/configs"
	"go.uber.org/zap"
	"time"
)

func NewDBConn(ctx context.Context, conf configs.DB, logger *zap.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.ConnectionURL())
	if err != nil {
		return nil, fmt.Errorf("[NewDBConn] sql open: %w", err)
	}

	if err = waitDBConn(ctx, conf, db, logger); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	return db, nil
}

const waitPingDBTimeout = time.Second * 40
const waitPingDBPolingInterval = time.Second * 2

func waitDBConn(ctx context.Context, conf configs.DB, db *sql.DB, logger *zap.Logger) error {
	ctx, cancel := context.WithTimeout(ctx, waitPingDBTimeout)
	defer cancel()

	conf.User = "user"
	conf.Password = "password"
	logger.Info("START wait DB connection", zap.String("database URL", conf.ConnectionURL()))

	attempts := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("[waitDBConn] ctx done, error: %w", ctx.Err())
		case <-time.After(waitPingDBPolingInterval):
			if err := db.PingContext(ctx); err != nil {
				attempts++
				logger.Info("wait DB connection", zap.Int("attempt", attempts))

				continue
			}

			return nil
		}
	}
}
