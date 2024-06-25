// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package api

type Delete struct {
	KamajiVersion     string `json:"kamaji_version"`
	KubernetesVersion string `json:"kubernetes_version"`
	TenantVersion     string `json:"tcp_version"`
	Status            string `json:"status"`
}
