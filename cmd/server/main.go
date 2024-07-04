// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/clastix/kamaji-telemetry/internal/ctx"
	"github.com/clastix/kamaji-telemetry/internal/db"
	"github.com/clastix/kamaji-telemetry/internal/webserver"
)

func main() {
	var args TelemetryArgs

	zapOpts, zapFs := &zap.Options{Development: true}, flag.NewFlagSet("zap", flag.ExitOnError)
	zapOpts.BindFlags(zapFs)
	logger := zap.New(zap.UseFlagOptions(zapOpts))

	var pool *pgxpool.Pool

	cmd := &cobra.Command{
		Use:   "kamaji-telemetry",
		Short: "HTTP web server storing Kamaji tenant Control Plane telemetry events",
		Long: "Storing Kamaji tenant Control Plane telemetry events such as creation, updates, and most used versions. " +
			"These info are used to collect anonymous data to be analyzed for the project development roadmap.",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, _ []string) (err error) {
			confConfig, confErr := pgxpool.ParseConfig(args.PsqlConnectionString)
			if confErr != nil {
				return errors.Wrap(err, "unable to parse connection string")
			}

			pool, err = pgxpool.NewWithConfig(cmd.Context(), confConfig)
			if err != nil {
				return errors.Wrap(err, "unable to create the pool")
			}

			pingCtx, cancel := ctx.GracefulTimeout(5 * time.Second)
			defer cancel()

			if err = pool.Ping(pingCtx); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}

			if initErr := db.CreateTableIfNotExists(cmd.Context(), pool); initErr != nil {
				return errors.Wrap(err, "unable to initialize table")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := webserver.New(cmd.Context(), logger, pool, args.Port); err != nil {
				return errors.Wrap(err, "failed to start server")
			}

			pool.Close()

			return nil
		},
	}

	cmd.Flags().IntVar(&args.Port, "listening-port", 8080, "Port the web server is listening on.")
	_ = cmd.MarkFlagRequired("listening-port")

	cmd.Flags().StringVar(&args.PsqlConnectionString, "psql-connection-string", "postgresql://user:password@postgres-db:5432/dbname", "PostgreSQL connection string.")
	_ = cmd.MarkFlagRequired("psql-connection-string")

	cmd.Flags().AddGoFlagSet(zapFs)

	if err := cmd.ExecuteContext(ctx.SignalNotified(context.Background())); err != nil {
		logger.Error(err, "error executing command")

		os.Exit(1)
	}

	logger.Info("shutdown complete")

	os.Exit(0)
}
