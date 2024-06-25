// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"github.com/clastix/kamaji-telemetry/api"
)

type noOp struct{}

func (n noOp) PushStats(context.Context, api.Stats) {
}

func (n noOp) PushCreate(context.Context, api.Create) {
}

func (n noOp) PushDelete(context.Context, api.Delete) {
}

func (n noOp) PushUpdate(context.Context, api.Update) {
}

func NewNewOp() Client {
	return &noOp{}
}
