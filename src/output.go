package main

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"
)

// generateOutputFiles generates JSON and HTML output files
func generateOutputFiles(config *ComparisonConfig) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write Cluster A data to file
	if err := writeJSONFile(fmt.Sprintf("%s/cluster-a-%s.json", config.OutputDir, config.ReportTimestamp), config.ClusterA.Data); err != nil {
		return fmt.Errorf("failed to write cluster-a.json: %w", err)
	}

	// Write Cluster B data to file
	if err := writeJSONFile(fmt.Sprintf("%s/cluster-b-%s.json", config.OutputDir, config.ReportTimestamp), config.ClusterB.Data); err != nil {
		return fmt.Errorf("failed to write cluster-b.json: %w", err)
	}

	// Generate HTML report
	if err := generateHTMLReport(config); err != nil {
		return fmt.Errorf("failed to generate HTML report: %w", err)
	}

	return nil
}

// writeJSONFile writes data to a JSON file
func writeJSONFile(filename string, data []map[string]interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

// generateHTMLReport creates an HTML report with embedded data
func generateHTMLReport(config *ComparisonConfig) error {
	// Create timestamp for filename
	filename := fmt.Sprintf("%s/k8s-comparison-report_%s.html", config.OutputDir, config.ReportTimestamp)

	// Convert data to JSON strings for embedding
	clusterAJSON, err := json.Marshal(config.ClusterA.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal cluster A data: %w", err)
	}

	clusterBJSON, err := json.Marshal(config.ClusterB.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal cluster B data: %w", err)
	}

	// Generate the HTML content
	htmlContent := generateHTMLTemplate(config, string(clusterAJSON), string(clusterBJSON), config.ReportTimestamp)

	// Write to file
	err = os.WriteFile(filename, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	fmt.Printf("📄 Generated HTML report: %s\n", filename)
	return nil
}

// generateResourceTags creates HTML tags for resources
func generateResourceTags(resources []string) string {
	var tags []string
	for _, resource := range resources {
		tags = append(tags, fmt.Sprintf(`<span class="resource-tag">%s</span>`, html.EscapeString(resource)))
	}
	return strings.Join(tags, "\n                        ")
}
