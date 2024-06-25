// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package ctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SignalNotified(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1)
	}()

	return ctx
}
