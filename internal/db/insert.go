// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

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

func Create(conn *pgx.Conn, payload []byte, logger logr.Logger) {
	ctx, cancel := kamajictx.GracefulTimeout(insertTimeout)
	defer cancel()

	if _, err := conn.Exec(ctx, insertStatement, uuid.New().String(), "create", payload); err != nil {
		logger.Error(err, "failed to insert create event")
	}
}

func Update(conn *pgx.Conn, payload []byte, logger logr.Logger) {
	ctx, cancel := kamajictx.GracefulTimeout(insertTimeout)
	defer cancel()

	if _, err := conn.Exec(ctx, insertStatement, uuid.New().String(), "update", payload); err != nil {
		logger.Error(err, "failed to insert update event")
	}
}

func Delete(conn *pgx.Conn, payload []byte, logger logr.Logger) {
	ctx, cancel := kamajictx.GracefulTimeout(insertTimeout)
	defer cancel()

	if _, err := conn.Exec(ctx, insertStatement, uuid.New().String(), "delete", payload); err != nil {
		logger.Error(err, "failed to insert delete event")
	}
}

func Stats(conn *pgx.Conn, payload []byte, logger logr.Logger, uuid string) {
	ctx, cancel := kamajictx.GracefulTimeout(insertTimeout)
	defer cancel()

	if _, err := conn.Exec(ctx, createOrUpdateStatement, uuid, "stats", payload); err != nil {
		logger.Error(err, "failed to insert delete event")
	}
}
