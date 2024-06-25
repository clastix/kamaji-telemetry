// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package ctx

import (
	"context"
	"time"
)

func GracefulTimeout(duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), duration)
}
