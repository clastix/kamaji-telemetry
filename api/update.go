// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package api

type Update struct {
	KamajiVersion     string `json:"kamaji_version"`
	KubernetesVersion string `json:"kubernetes_version"`
	OldTenantVersion  string `json:"old_tcp_version"`
	NewTenantVersion  string `json:"new_tcp_version"`
}
