// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/clastix/kamaji-telemetry/api"
)

const (
	ContentTypeJSON = "application/json"
)

func New(baseClient http.Client, baseURL string) Client {
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &HTTP{
		baseClient: baseClient,
		baseURL:    baseURL,
	}
}

type HTTP struct {
	baseClient http.Client
	baseURL    string
}

func (h *HTTP) doRequest(ctx context.Context, url string, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.baseURL+"/"+url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", ContentTypeJSON)
	req.Header.Set("User-Agent", "clastix/kamaji-telemetry")

	res, err := h.baseClient.Do(req)
	if err != nil {
		return
	}

	_ = res.Body.Close()
}

func (h *HTTP) PushStats(ctx context.Context, payload api.Stats) {
	h.doRequest(ctx, "stats", payload)
}

func (h *HTTP) PushCreate(ctx context.Context, payload api.Create) {
	h.doRequest(ctx, "create", payload)
}

func (h *HTTP) PushDelete(ctx context.Context, payload api.Delete) {
	h.doRequest(ctx, "delete", payload)
}

func (h *HTTP) PushUpdate(ctx context.Context, payload api.Update) {
	h.doRequest(ctx, "update", payload)
}
