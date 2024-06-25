// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func CreateTableIfNotExists(ctx context.Context, c *pgx.Conn) error {
	if _, err := c.Exec(ctx, CreateStatement); err != nil {
		return errors.Wrap(err, "cannot create or update required table")
	}

	return nil
}
