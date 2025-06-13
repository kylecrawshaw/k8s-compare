package main

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Describe("ClusterConfig struct", func() {
		Context("when creating a new ClusterConfig", func() {
			It("should initialize with correct fields", func() {
				config := ClusterConfig{
					Context:    "test-context",
					Namespaces: []string{"test-namespace"},
					Resources:  []string{"pods", "services"},
					Data:       []map[string]interface{}{{"key": "value"}},
				}

				Expect(config.Context).To(Equal("test-context"))
				Expect(config.Namespaces).To(Equal([]string{"test-namespace"}))
				Expect(config.Resources).To(Equal([]string{"pods", "services"}))
				Expect(config.Data).To(HaveLen(1))
			})

			It("should allow empty values", func() {
				config := ClusterConfig{}

				Expect(config.Context).To(Equal(""))
				Expect(config.Namespaces).To(BeNil())
				Expect(config.Resources).To(BeNil())
				Expect(config.Data).To(BeNil())
			})
		})

		Context("when modifying ClusterConfig fields", func() {
			It("should allow field updates", func() {
				config := ClusterConfig{
					Context:    "original-context",
					Namespaces: []string{"original-namespace"},
					Resources:  []string{"pods"},
					Data:       []map[string]interface{}{},
				}

				config.Context = "updated-context"
				config.Namespaces = []string{"updated-namespace"}
				config.Resources = []string{"services"}
				config.Data = []map[string]interface{}{{"updated": "data"}}

				Expect(config.Context).To(Equal("updated-context"))
				Expect(config.Namespaces).To(Equal([]string{"updated-namespace"}))
				Expect(config.Resources).To(Equal([]string{"services"}))
				Expect(config.Data).To(HaveLen(1))
			})
		})
	})

	Describe("ComparisonConfig struct", func() {
		Context("when creating a new ComparisonConfig", func() {
			It("should initialize with correct fields", func() {
				clusterA := ClusterConfig{
					Context:    "context1",
					Namespaces: []string{"namespace1"},
					Resources:  []string{"pods"},
				}
				clusterB := ClusterConfig{
					Context:    "context2",
					Namespaces: []string{"namespace2"},
					Resources:  []string{"services"},
				}

				config := ComparisonConfig{
					ClusterA: clusterA,
					ClusterB: clusterB,
				}

				Expect(config.ClusterA).To(Equal(clusterA))
				Expect(config.ClusterB).To(Equal(clusterB))
			})

			It("should allow empty clusters", func() {
				config := ComparisonConfig{}

				Expect(config.ClusterA).To(Equal(ClusterConfig{}))
				Expect(config.ClusterB).To(Equal(ClusterConfig{}))
			})
		})

		Context("when working with cluster configurations", func() {
			It("should handle different cluster contexts", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{
						Context:    "gke_project_zone_cluster-a",
						Namespaces: []string{"default", "kube-system"},
						Resources:  []string{"pods", "services"},
					},
					ClusterB: ClusterConfig{
						Context:    "minikube",
						Namespaces: []string{"default"},
						Resources:  []string{"pods", "deployments"},
					},
				}

				Expect(config.ClusterA.Context).To(Equal("gke_project_zone_cluster-a"))
				Expect(config.ClusterB.Context).To(Equal("minikube"))
				Expect(config.ClusterA.Namespaces).To(HaveLen(2))
				Expect(config.ClusterB.Namespaces).To(HaveLen(1))
			})

			It("should handle empty namespaces and resources", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "cluster-a"},
					ClusterB: ClusterConfig{Context: "cluster-b"},
				}

				Expect(config.ClusterA.Namespaces).To(BeNil())
				Expect(config.ClusterA.Resources).To(BeNil())
				Expect(config.ClusterB.Namespaces).To(BeNil())
				Expect(config.ClusterB.Resources).To(BeNil())
			})
		})

		Context("when comparing ComparisonConfig instances", func() {
			It("should be equal when all fields match", func() {
				clusterA := ClusterConfig{Context: "ctx1", Namespaces: []string{"ns1"}}
				clusterB := ClusterConfig{Context: "ctx2", Namespaces: []string{"ns2"}}

				config1 := ComparisonConfig{
					ClusterA: clusterA,
					ClusterB: clusterB,
				}

				config2 := ComparisonConfig{
					ClusterA: clusterA,
					ClusterB: clusterB,
				}

				Expect(config1).To(Equal(config2))
			})

			It("should not be equal when clusters differ", func() {
				clusterA := ClusterConfig{Context: "ctx1", Namespaces: []string{"ns1"}}
				clusterB := ClusterConfig{Context: "ctx2", Namespaces: []string{"ns2"}}
				clusterC := ClusterConfig{Context: "ctx3", Namespaces: []string{"ns3"}}

				config1 := ComparisonConfig{
					ClusterA: clusterA,
					ClusterB: clusterB,
				}

				config2 := ComparisonConfig{
					ClusterA: clusterA,
					ClusterB: clusterC,
				}

				Expect(config1).NotTo(Equal(config2))
			})

			It("should not be equal when cluster contexts differ", func() {
				config1 := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "ctx1"},
					ClusterB: ClusterConfig{Context: "ctx2"},
				}

				config2 := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "ctx1"},
					ClusterB: ClusterConfig{Context: "ctx3"},
				}

				Expect(config1).NotTo(Equal(config2))
			})
		})
	})
})
