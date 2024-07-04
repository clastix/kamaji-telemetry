// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	CreateStatement = `CREATE TABLE IF NOT EXISTS events (
  uuid varchar(36) NOT NULL,
  event_type varchar NOT NULL,
  payload jsonb not null,
  created_at timestamp with time zone NOT NULL,
  PRIMARY KEY (uuid)
)`
)

func CreateTableIfNotExists(ctx context.Context, p *pgxpool.Pool) error {
	err := p.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		tx, txErr := conn.Begin(ctx)
		if txErr != nil {
			return errors.Wrap(txErr, "transaction failed")
		}

		if _, err := tx.Exec(ctx, CreateStatement); err != nil {
			_ = tx.Rollback(ctx)

			return errors.Wrap(err, "cannot create or update required table")
		}

		return tx.Commit(ctx)
	})
	if err != nil {
		return errors.Wrap(err, "cannot create table")
	}

	return nil
}
