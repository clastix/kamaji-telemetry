// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package api

const (
	NetworkProfileIngress   = "ingress"
	NetworkProfileLB        = "load_balancer"
	NetworkProfileNodePort  = "node_port"
	NetworkProfileClusterIP = "cluster_ip"
)

type Create struct {
	KamajiVersion     string `json:"kamaji_version"`
	KubernetesVersion string `json:"kubernetes_version"`
	TenantVersion     string `json:"tcp_version"`
	ClusterAPIOwned   bool   `json:"cluster_api"`
	NetworkProfile    string `json:"network_profile"`
}
