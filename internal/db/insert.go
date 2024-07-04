// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	kamajictx "github.com/clastix/kamaji-telemetry/internal/ctx"
)

const (
	insertStatement         = `insert into events (uuid, event_type, payload, created_at) values ($1, $2, $3, now())`
	insertTimeout           = 5 * time.Second
	createOrUpdateStatement = `INSERT INTO events (uuid, event_type, payload, created_at) 
VALUES ($1, $2, $3, now())
ON CONFLICT (uuid) DO UPDATE
SET payload = $3, created_at = now()`
)

func Create(conn *pgxpool.Pool, payload []byte, logger logr.Logger) {
	err := execWithTransaction(conn, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, insertStatement, uuid.New().String(), "create", payload)

		return err //nolint:wrapcheck
	})
	if err != nil {
		logger.Error(err, "failed to insert create event")
	}
}

func Update(conn *pgxpool.Pool, payload []byte, logger logr.Logger) {
	err := execWithTransaction(conn, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, insertStatement, uuid.New().String(), "update", payload)

		return err //nolint:wrapcheck
	})
	if err != nil {
		logger.Error(err, "failed to insert update event")
	}
}

func Delete(conn *pgxpool.Pool, payload []byte, logger logr.Logger) {
	err := execWithTransaction(conn, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, insertStatement, uuid.New().String(), "delete", payload)

		return err //nolint:wrapcheck
	})
	if err != nil {
		logger.Error(err, "failed to insert delete event")
	}
}

func Stats(conn *pgxpool.Pool, payload []byte, logger logr.Logger, uuid string) {
	err := execWithTransaction(conn, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, createOrUpdateStatement, uuid, "stats", payload)

		return err //nolint:wrapcheck
	})
	if err != nil {
		logger.Error(err, "failed to insert stats event")
	}
}

func execWithTransaction(pool *pgxpool.Pool, fn func(context.Context, pgx.Tx) error) error {
	ctx, cancel := kamajictx.GracefulTimeout(insertTimeout)
	defer cancel()

	tx, txErr := pool.Begin(ctx)
	if txErr != nil {
		return errors.Wrap(txErr, "failed to begin transaction for delete")
	}

	if err := fn(ctx, tx); err != nil {
		_ = tx.Rollback(context.Background())

		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}
