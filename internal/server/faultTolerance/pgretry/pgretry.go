package pgretry

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

const (
	maxRetries          = 3
	retryInterval       = 5
	maxBackoffInterval  = 10 * time.Second
	randomizationFactor = 0
	multiplier          = 1
)

type PgxRetry struct {
	dbpool   *pgxpool.Pool
	log      *zap.SugaredLogger
	duration []int
}

func NewPgxRetry(conn *pgxpool.Pool, log *zap.SugaredLogger) *PgxRetry {
	if conn == nil {
		return nil
	}

	return &PgxRetry{
		dbpool: conn,
		log:    log,
	}
}

func (pgr *PgxRetry) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()
	operation := func() (pgconn.CommandTag, error) {
		tag, err := pgr.dbpool.Exec(ctx, sql, arguments...)

		if err != nil {
			end := time.Now()
			duration := end.Sub(start)
			var connectErr *pgconn.ConnectError
			if errors.As(err, &connectErr) {
				log.Debugf("ошибка подключения к базе; Пробуем еще раз, прошло времени: %v сек", duration.Seconds())

				return tag, backoff.RetryAfter(retryInterval)
			}

			log.Debugf("ошибка фатальна; Прошло времени: %v сек", duration.Seconds())
			return tag, backoff.Permanent(err)
		}

		return tag, nil
	}

	return backoff.Retry(ctx, operation, backoff.WithBackOff(pgr.getBackOffOptions()), backoff.WithMaxTries(maxRetries))
}

func (pgr *PgxRetry) QueryRow(ctx context.Context, sql string, args ...interface{}) (pgx.Row, error) {
	operation := func() (pgx.Row, error) {
		row := pgr.dbpool.QueryRow(ctx, sql, args...)

		return row, nil
	}

	return backoff.Retry(ctx, operation, backoff.WithBackOff(pgr.getBackOffOptions()), backoff.WithMaxTries(maxRetries))
}

func (pgr *PgxRetry) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()
	operation := func() (pgx.Rows, error) {
		rows, err := pgr.dbpool.Query(ctx, sql, args...)

		if err != nil {
			end := time.Now()
			duration := end.Sub(start)
			var connectErr *pgconn.ConnectError
			if errors.As(err, &connectErr) {
				log.Debugf("ошибка подключения к базе; Пробуем еще раз, прошло времени: %v сек", duration.Seconds())

				return rows, backoff.RetryAfter(retryInterval)
			}

			log.Debugf("ошибка фатальна; Прошло времени: %v сек", duration.Seconds())
			return rows, backoff.Permanent(err)
		}

		return rows, nil
	}

	return backoff.Retry(ctx, operation, backoff.WithBackOff(pgr.getBackOffOptions()), backoff.WithMaxTries(maxRetries))
}

func (pgr *PgxRetry) Begin(ctx context.Context) (pgx.Tx, error) {
	start := time.Now()
	operation := func() (pgx.Tx, error) {
		tx, err := pgr.dbpool.Begin(ctx)

		if err != nil {
			end := time.Now()
			duration := end.Sub(start)
			var connectErr *pgconn.ConnectError
			if errors.As(err, &connectErr) {
				log.Debugf("ошибка подключения к базе; Пробуем еще раз, прошло времени: %v сек", duration.Seconds())

				return tx, backoff.RetryAfter(retryInterval)
			}

			log.Debugf("ошибка фатальна; Прошло времени: %v сек", duration.Seconds())
			return tx, backoff.Permanent(err)
		}

		return tx, nil
	}

	return backoff.Retry(ctx, operation, backoff.WithBackOff(pgr.getBackOffOptions()), backoff.WithMaxTries(maxRetries))
}

func (pgr *PgxRetry) Ping(ctx context.Context) error {
	start := time.Now()
	operation := func() (bool, error) {
		err := pgr.dbpool.Ping(ctx)

		if err != nil {
			end := time.Now()
			duration := end.Sub(start)
			var connectErr *pgconn.ConnectError
			if errors.As(err, &connectErr) {
				log.Debugf("ошибка подключения к базе; Пробуем еще раз, прошло времени: %v сек", duration.Seconds())

				return false, backoff.RetryAfter(retryInterval)
			}

			log.Debugf("ошибка фатальна; Прошло времени: %v сек", duration.Seconds())
			return false, backoff.Permanent(err)
		}

		return true, nil
	}

	_, err := backoff.Retry(ctx, operation, backoff.WithBackOff(pgr.getBackOffOptions()), backoff.WithMaxTries(maxRetries))
	return err
}

func (pgr *PgxRetry) getBackOffOptions() *backoff.ExponentialBackOff {
	return &backoff.ExponentialBackOff{
		InitialInterval:     retryInterval * time.Second,
		RandomizationFactor: randomizationFactor,
		Multiplier:          multiplier,
		MaxInterval:         maxBackoffInterval,
	}
}
