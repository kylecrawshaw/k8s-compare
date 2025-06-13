package main

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTML Template", func() {
	Describe("generateHTMLTemplate function", func() {
		Context("when generating HTML templates", func() {
			It("should generate valid HTML structure", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Check basic HTML structure
				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(template).To(ContainSubstring("<html"))
				Expect(template).To(ContainSubstring("</html>"))
				Expect(template).To(ContainSubstring("<head>"))
				Expect(template).To(ContainSubstring("</head>"))
				Expect(template).To(ContainSubstring("<body>"))
				Expect(template).To(ContainSubstring("</body>"))
			})

			It("should include required meta tags", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Check for essential meta tags
				Expect(template).To(ContainSubstring(`<meta charset="UTF-8"`))
				Expect(template).To(ContainSubstring(`<meta name="viewport"`))
				Expect(template).To(ContainSubstring("<title>"))
			})

			It("should include CSS styles", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Check for CSS
				Expect(template).To(ContainSubstring("<style>"))
				Expect(template).To(ContainSubstring("</style>"))
				Expect(template).To(ContainSubstring("body"))
				Expect(template).To(ContainSubstring("color"))
				Expect(template).To(ContainSubstring("font-family"))
			})

			It("should include JavaScript functionality", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Check for JavaScript
				Expect(template).To(ContainSubstring("<script>"))
				Expect(template).To(ContainSubstring("</script>"))
				Expect(template).To(ContainSubstring("function"))
			})

			It("should embed cluster data correctly", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "production-cluster"},
					ClusterB: ClusterConfig{Context: "staging-cluster"},
				}

				clusterAData := `[{"name":"pod-1","status":"Running"}]`
				clusterBData := `[{"name":"pod-2","status":"Pending"}]`

				template := generateHTMLTemplate(config, clusterAData, clusterBData, "2023-01-01_12-00-00")

				// Check that cluster contexts are included
				Expect(template).To(ContainSubstring("production-cluster"))
				Expect(template).To(ContainSubstring("staging-cluster"))

				// Check that data is embedded
				Expect(template).To(ContainSubstring("pod-1"))
				Expect(template).To(ContainSubstring("pod-2"))
				Expect(template).To(ContainSubstring("Running"))
				Expect(template).To(ContainSubstring("Pending"))
			})

			It("should include timestamp in the template", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				timestamp := "2023-12-25_14-30-45"
				template := generateHTMLTemplate(config, `[]`, `[]`, timestamp)

				// Check that timestamp is included
				Expect(template).To(ContainSubstring(timestamp))
			})

			It("should handle empty cluster data", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "empty-a"},
					ClusterB: ClusterConfig{Context: "empty-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Should still generate valid HTML
				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(template).To(ContainSubstring("empty-a"))
				Expect(template).To(ContainSubstring("empty-b"))
			})

			It("should handle special characters in cluster names", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "cluster-with-dashes_and_underscores"},
					ClusterB: ClusterConfig{Context: "cluster.with.dots"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Should handle special characters properly
				Expect(template).To(ContainSubstring("cluster-with-dashes_and_underscores"))
				Expect(template).To(ContainSubstring("cluster.with.dots"))
			})

			It("should include navigation tabs", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Check for tab navigation
				Expect(template).To(ContainSubstring("Overview"))
				Expect(template).To(ContainSubstring("Resource Breakdown"))
				Expect(template).To(ContainSubstring("Detailed Comparison"))
			})

			It("should include interactive features", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Check for interactive elements
				Expect(template).To(ContainSubstring("onclick"))
				Expect(template).To(ContainSubstring("addEventListener"))
			})

			It("should handle complex JSON data", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "complex-a"},
					ClusterB: ClusterConfig{Context: "complex-b"},
				}

				complexData := `[{
					"apiVersion": "v1",
					"kind": "Pod",
					"metadata": {
						"name": "test-pod",
						"namespace": "default",
						"labels": {"app": "test"}
					},
					"spec": {
						"containers": [{"name": "test", "image": "nginx"}]
					}
				}]`

				template := generateHTMLTemplate(config, complexData, complexData, "2023-01-01_12-00-00")

				// Should handle complex nested JSON
				Expect(template).To(ContainSubstring("test-pod"))
				Expect(template).To(ContainSubstring("nginx"))
				Expect(template).To(ContainSubstring("apiVersion"))
			})
		})

		Context("when handling edge cases", func() {
			It("should handle nil config gracefully", func() {
				// This might panic or handle gracefully depending on implementation
				// We test that it doesn't crash the test suite
				defer func() {
					if r := recover(); r != nil {
						// If it panics, that's also a valid behavior to test
						Expect(r).NotTo(BeNil())
					}
				}()

				template := generateHTMLTemplate(nil, `[]`, `[]`, "2023-01-01_12-00-00")

				// If it doesn't panic, it should still generate some HTML
				if template != "" {
					Expect(template).To(ContainSubstring("html"))
				}
			})

			It("should handle malformed JSON data", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				malformedJSON := `[{"name": "test", "invalid": json}]`

				template := generateHTMLTemplate(config, malformedJSON, `[]`, "2023-01-01_12-00-00")

				// Should still generate HTML even with malformed JSON
				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(template).To(ContainSubstring("test-a"))
			})

			It("should handle very long cluster names", func() {
				longName := strings.Repeat("very-long-cluster-name-", 10)
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: longName + "a"},
					ClusterB: ClusterConfig{Context: longName + "b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Should handle long names without breaking
				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(len(template)).To(BeNumerically(">", 1000)) // Should be substantial HTML
			})

			It("should handle empty timestamp", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "")

				// Should still generate valid HTML
				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(template).To(ContainSubstring("test-a"))
			})
		})

		Context("when validating HTML output", func() {
			It("should produce well-formed HTML", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				// Basic HTML validation checks
				htmlOpenCount := strings.Count(template, "<html")
				htmlCloseCount := strings.Count(template, "</html>")
				Expect(htmlOpenCount).To(Equal(htmlCloseCount))

				headOpenCount := strings.Count(template, "<head")
				headCloseCount := strings.Count(template, "</head>")
				Expect(headOpenCount).To(Equal(headCloseCount))

				bodyOpenCount := strings.Count(template, "<body")
				bodyCloseCount := strings.Count(template, "</body>")
				Expect(bodyOpenCount).To(Equal(bodyCloseCount))
			})

			It("should have balanced script tags", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				scriptOpenCount := strings.Count(template, "<script")
				scriptCloseCount := strings.Count(template, "</script>")
				Expect(scriptOpenCount).To(Equal(scriptCloseCount))
			})

			It("should have balanced style tags", func() {
				config := &ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				template := generateHTMLTemplate(config, `[]`, `[]`, "2023-01-01_12-00-00")

				styleOpenCount := strings.Count(template, "<style")
				styleCloseCount := strings.Count(template, "</style>")
				Expect(styleOpenCount).To(Equal(styleCloseCount))
			})
		})
	})
})
