// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"github.com/clastix/kamaji-telemetry/api"
)

type Client interface {
	PushStats(ctx context.Context, payload api.Stats)
	PushCreate(ctx context.Context, payload api.Create)
	PushDelete(ctx context.Context, payload api.Delete)
	PushUpdate(ctx context.Context, payload api.Update)
}
