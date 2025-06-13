package main

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Output", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "k8s-compare-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("generateResourceTags function", func() {
		Context("when generating tags for resources", func() {
			It("should generate correct tags for pods", func() {
				resources := []string{"pods", "services", "deployments"}
				tags := generateResourceTags(resources)

				Expect(tags).To(ContainSubstring("pods"))
				Expect(tags).To(ContainSubstring("services"))
				Expect(tags).To(ContainSubstring("deployments"))
				Expect(tags).To(ContainSubstring("resource-tag"))
			})

			It("should handle empty resource list", func() {
				resources := []string{}
				tags := generateResourceTags(resources)

				Expect(tags).To(Equal(""))
			})

			It("should handle single resource", func() {
				resources := []string{"pods"}
				tags := generateResourceTags(resources)

				Expect(tags).To(ContainSubstring("pods"))
				Expect(tags).To(ContainSubstring("resource-tag"))
			})

			It("should handle special characters in resource names", func() {
				resources := []string{"custom-resources.v1", "api-extensions"}
				tags := generateResourceTags(resources)

				Expect(tags).To(ContainSubstring("custom-resources.v1"))
				Expect(tags).To(ContainSubstring("api-extensions"))
			})
		})
	})

	Describe("generateOutputFiles function", func() {
		Context("when generating output files", func() {
			It("should create JSON files successfully", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{
						Context:    "test-cluster-a",
						Namespaces: []string{"default"},
						Resources:  []string{"pods"},
						Data:       []map[string]interface{}{{"name": "pod-a", "status": "Running"}},
					},
					ClusterB: ClusterConfig{
						Context:    "test-cluster-b",
						Namespaces: []string{"default"},
						Resources:  []string{"pods"},
						Data:       []map[string]interface{}{{"name": "pod-b", "status": "Pending"}},
					},
				}

				// Change to temp directory
				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)
				os.Chdir(tempDir)

				err := generateOutputFiles(&config)
				Expect(err).NotTo(HaveOccurred())

				// Check if JSON files were created
				clusterAFile := filepath.Join(tempDir, "cluster-a.json")
				Expect(clusterAFile).To(BeAnExistingFile())

				clusterBFile := filepath.Join(tempDir, "cluster-b.json")
				Expect(clusterBFile).To(BeAnExistingFile())
			})

			It("should handle empty config", func() {
				config := ComparisonConfig{}

				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)
				os.Chdir(tempDir)

				err := generateOutputFiles(&config)
				Expect(err).NotTo(HaveOccurred())

				clusterAFile := filepath.Join(tempDir, "cluster-a.json")
				Expect(clusterAFile).To(BeAnExistingFile())
			})
		})
	})

	Describe("generateHTMLReport function", func() {
		Context("when generating HTML report", func() {
			It("should create HTML file with correct content", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{
						Context:    "gke_project_zone_cluster-a",
						Namespaces: []string{"default", "kube-system"},
						Resources:  []string{"pods", "services"},
						Data: []map[string]interface{}{
							{"name": "pod-a", "namespace": "default", "status": "Running"},
							{"name": "service-a", "namespace": "default", "type": "ClusterIP"},
						},
					},
					ClusterB: ClusterConfig{
						Context:    "minikube",
						Namespaces: []string{"default"},
						Resources:  []string{"pods", "services"},
						Data: []map[string]interface{}{
							{"name": "pod-b", "namespace": "default", "status": "Pending"},
						},
					},
				}

				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)
				os.Chdir(tempDir)

				err := generateHTMLReport(&config)
				Expect(err).NotTo(HaveOccurred())

				// Check that an HTML file was created (filename includes timestamp)
				files, err := os.ReadDir(tempDir)
				Expect(err).NotTo(HaveOccurred())

				var htmlFiles []string
				for _, file := range files {
					if filepath.Ext(file.Name()) == ".html" {
						htmlFiles = append(htmlFiles, file.Name())
					}
				}

				Expect(htmlFiles).To(HaveLen(1))

				// Read and verify HTML content
				htmlFile := filepath.Join(tempDir, htmlFiles[0])
				content, err := os.ReadFile(htmlFile)
				Expect(err).NotTo(HaveOccurred())

				htmlContent := string(content)
				Expect(htmlContent).To(ContainSubstring("Kubernetes Resource Comparison"))
				Expect(htmlContent).To(ContainSubstring("gke_project_zone_cluster-a"))
				Expect(htmlContent).To(ContainSubstring("minikube"))
			})

			It("should handle config with no data", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "cluster-a"},
					ClusterB: ClusterConfig{Context: "cluster-b"},
				}

				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)
				os.Chdir(tempDir)

				err := generateHTMLReport(&config)
				Expect(err).NotTo(HaveOccurred())

				// Check that an HTML file was created
				files, err := os.ReadDir(tempDir)
				Expect(err).NotTo(HaveOccurred())

				var htmlFiles []string
				for _, file := range files {
					if filepath.Ext(file.Name()) == ".html" {
						htmlFiles = append(htmlFiles, file.Name())
					}
				}

				Expect(htmlFiles).To(HaveLen(1))
			})
		})
	})

	Describe("generateHTMLTemplate function", func() {
		Context("when generating HTML template", func() {
			It("should return valid HTML template", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{
						Context:   "test-cluster-a",
						Resources: []string{"pods", "services"},
					},
					ClusterB: ClusterConfig{
						Context:   "test-cluster-b",
						Resources: []string{"pods", "services"},
					},
				}

				template := generateHTMLTemplate(&config, `[{"name":"test"}]`, `[{"name":"test2"}]`, "2023-01-01_12-00-00")

				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(template).To(ContainSubstring("<html"))
				Expect(template).To(ContainSubstring("</html>"))
				Expect(template).To(ContainSubstring("Kubernetes Resource Comparison"))
				Expect(template).To(ContainSubstring("test-cluster-a"))
				Expect(template).To(ContainSubstring("test-cluster-b"))
			})

			It("should include JavaScript functionality", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "cluster-a"},
					ClusterB: ClusterConfig{Context: "cluster-b"},
				}

				template := generateHTMLTemplate(&config, "[]", "[]", "2023-01-01_12-00-00")

				Expect(template).To(ContainSubstring("<script>"))
				Expect(template).To(ContainSubstring("</script>"))
				Expect(template).To(ContainSubstring("function"))
			})

			It("should include CSS styling", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "cluster-a"},
					ClusterB: ClusterConfig{Context: "cluster-b"},
				}

				template := generateHTMLTemplate(&config, "[]", "[]", "2023-01-01_12-00-00")

				Expect(template).To(ContainSubstring("<style>"))
				Expect(template).To(ContainSubstring("</style>"))
				Expect(template).To(ContainSubstring("body"))
				Expect(template).To(ContainSubstring("color"))
			})

			It("should handle empty config gracefully", func() {
				config := ComparisonConfig{}

				template := generateHTMLTemplate(&config, "[]", "[]", "2023-01-01_12-00-00")

				Expect(template).To(ContainSubstring("<!DOCTYPE html>"))
				Expect(template).To(ContainSubstring("Kubernetes Resource Comparison"))
			})
		})
	})

	Describe("file operations", func() {
		Context("when working with file paths", func() {
			It("should create files in current directory by default", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)
				os.Chdir(tempDir)

				err := generateOutputFiles(&config)
				Expect(err).NotTo(HaveOccurred())

				files, err := os.ReadDir(tempDir)
				Expect(err).NotTo(HaveOccurred())

				var fileNames []string
				for _, file := range files {
					fileNames = append(fileNames, file.Name())
				}

				Expect(fileNames).To(ContainElement("cluster-a.json"))
				Expect(fileNames).To(ContainElement("cluster-b.json"))
			})
		})

		Context("when handling timestamps", func() {
			It("should include timestamp in generated HTML files", func() {
				config := ComparisonConfig{
					ClusterA: ClusterConfig{Context: "test-a"},
					ClusterB: ClusterConfig{Context: "test-b"},
				}

				originalDir, _ := os.Getwd()
				defer os.Chdir(originalDir)
				os.Chdir(tempDir)

				beforeTime := time.Now()
				err := generateOutputFiles(&config)
				Expect(err).NotTo(HaveOccurred())

				// Check that HTML file includes timestamp in filename
				files, err := os.ReadDir(tempDir)
				Expect(err).NotTo(HaveOccurred())

				var htmlFiles []string
				for _, file := range files {
					if filepath.Ext(file.Name()) == ".html" {
						htmlFiles = append(htmlFiles, file.Name())
					}
				}

				Expect(htmlFiles).To(HaveLen(1))
				// The filename should contain the year at minimum
				Expect(htmlFiles[0]).To(ContainSubstring(beforeTime.Format("2006")))
			})
		})
	})
})
