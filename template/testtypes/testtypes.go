// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package testtypes

//go:generate cue get go --local github.com/errordeveloper/cue-utils/template/testtypes

type Cluster struct {
	Metadata ClusterMeta `json:"metadata"`
	Spec     ClusterSpec `json:"spec"`
}
type ClusterMeta struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
type ClusterSpec struct {
	Location   string  `json:"location"`
	SubnetCIDR *string `json:"subnetCIDR,omitempty"`
}
