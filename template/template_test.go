// Copyright 2020 Authors of Cilium
// Copyright 2022 Ilya Dmitrichenko
// SPDX-License-Identifier: Apache-2.0

package template_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	. "github.com/errordeveloper/cue-utils/template"
	"github.com/errordeveloper/cue-utils/template/testtypes"
)

func TestGenerator(t *testing.T) {
	g := NewGomegaWithT(t)

	primaryGen := NewGenerator("./testassets")

	err := primaryGen.CompileAndValidate()
	g.Expect(err).To(Not(HaveOccurred()))

	cidr := "10.128.0.0/20"
	baseGen, err := primaryGen.WithDefaults(&testtypes.Cluster{
		Spec: testtypes.ClusterSpec{
			SubnetCIDR: &cidr,
		},
	})
	g.Expect(err).To(Not(HaveOccurred()))

	{
		cluster := testtypes.Cluster{}
		cluster.Metadata.Name = "foo1"
		cluster.Metadata.Namespace = "default"
		cluster.Spec.Location = "us-central1-a"

		gen, err := baseGen.WithResource(cluster)
		g.Expect(err).To(Not(HaveOccurred()))

		js, err := gen.RenderJSON()
		g.Expect(err).To(Not(HaveOccurred()))

		// default CIDR will be used
		g.Expect(js).To(MatchJSON(expectedWithCIDR("10.128.0.0/20")))
	}

	{
		cluster := testtypes.Cluster{}
		cluster.Metadata.Name = "foo1"
		cluster.Metadata.Namespace = "default"
		cluster.Spec.Location = "us-central1-a"
		cluster.Spec.SubnetCIDR = new(string)
		*cluster.Spec.SubnetCIDR = "10.128.0.0/16"

		gen, err := baseGen.WithResource(cluster)
		g.Expect(err).To(Not(HaveOccurred()))

		js, err := gen.RenderJSON()
		g.Expect(err).To(Not(HaveOccurred()))

		// given CIDR will override the default
		g.Expect(js).To(MatchJSON(expectedWithCIDR("10.128.0.0/16")))
	}

	{
		cluster := map[string]string{
			"foo": "bar",
		}

		_, err := baseGen.WithResource(cluster)
		g.Expect(err).To(HaveOccurred())

		g.Expect(err.Error()).To(Equal(`error after filling "resource": resource: field not allowed: foo`))

	}

	{
		cluster := map[string]string{}

		gen, err := baseGen.WithResource(cluster)
		g.Expect(err).To(Not(HaveOccurred()))

		_, err = gen.RenderJSON()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`cue: marshal error: template.items.0.metadata.namespace: invalid interpolation: non-concrete value string (type string)`))
	}

	{

		_, err := baseGen.WithResource(0)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal(`error after filling "resource": resource: invalid interpolation: conflicting values 0 and {metadata:#ClusterMeta,spec:#ClusterSpec} (mismatched types int and struct) (and 3 more errors)`))
	}

	{
		gen := NewGenerator("")

		err := gen.CompileAndValidate()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal("failed to load instances (dir: \"\", args: [.]): no CUE files in .\n"))
	}

	{
		gen := NewGenerator("./nonexistent")

		err := gen.CompileAndValidate()

		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(Equal("failed to load instances (dir: \"./nonexistent\", args: [.]): cannot find package \".\"\n"))
	}
}

func expectedWithCIDR(cidr string) string {
	const jsfmt = `
	{
		"kind": "List",
		"apiVersion": "v1",
		"items": [
			{
				"metadata": {
					"name": "foo1",
					"namespace": "default",
					"labels": {
						"cluster": "foo1"
					},
					"annotations": {
						"cnrm.cloud.google.com/remove-default-node-pool": "false"
					}
				},
				"spec": {
					"location": "us-central1-a",
					"networkRef": {
						"name": "foo1"
					},
					"subnetworkRef": {
						"name": "foo1"
					},
					"initialNodeCount": 1,
					"loggingService": "logging.googleapis.com/kubernetes",
					"monitoringService": "monitoring.googleapis.com/kubernetes",
					"masterAuth": {
						"clientCertificateConfig": {
							"issueClientCertificate": false
						}
					}
				},
				"kind": "ContainerCluster",
				"apiVersion": "container.cnrm.cloud.google.com/v1beta1"
			},
			{
				"metadata": {
					"name": "foo1",
					"namespace": "default",
					"labels": {
						"cluster": "foo1"
					}
				},
				"spec": {
					"routingMode": "REGIONAL",
					"autoCreateSubnetworks": false,
					"deleteDefaultRoutesOnCreate": false
				},
				"kind": "ComputeNetwork",
				"apiVersion": "compute.cnrm.cloud.google.com/v1beta1"
			},
			{
				"metadata": {
					"name": "foo1",
					"namespace": "default",
					"labels": {
						"cluster": "foo1"
					}
				},
				"spec": {
					"networkRef": {
						"name": "foo1"
					},
					"region": "us-central1",
					"ipCidrRange": "%s"
				},
				"kind": "ComputeSubnetwork",
				"apiVersion": "compute.cnrm.cloud.google.com/v1beta1"
			}
		]
	}`
	return fmt.Sprintf(jsfmt, cidr)
}

func TestGeneratorWithImportPath(t *testing.T) {
	g := NewGomegaWithT(t)
	g.Expect(NewGenerator("./", "github.com/errordeveloper/cue-utils/template/testtypes").CompileAndValidate()).To(Succeed())
}
