// Copyright 2024 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package api

type Stats struct {
	UUID                string             `json:"uuid"`
	KamajiVersion       string             `json:"kamaji_version"`
	KubernetesVersion   string             `json:"kubernetes_version"`
	TenantControlPlanes TenantControlPlane `json:"tenant_control_plane"`
	Datastores          Datastores         `json:"datastores"`
}

type TenantControlPlane struct {
	Running   uint64 `json:"running"`
	Upgrading uint64 `json:"upgrading"`
	NotReady  uint64 `json:"not_ready"`
	Sleeping  uint64 `json:"sleeping"`
}

type Datastores struct {
	Etcd       Datastore `json:"etcd"`
	PostgreSQL Datastore `json:"postgres"`
	MySQL      Datastore `json:"mysql"`
	NATS       Datastore `json:"nats"`
}

type Datastore struct {
	ShardCount  int `json:"count"`
	UsedByCount int `json:"used_by_count"`
}
