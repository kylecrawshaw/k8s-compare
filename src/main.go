package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "k8s-compare",
		Short: "Compare resources between two Kubernetes clusters",
		Long:  `A CLI tool to interactively select contexts, namespaces, and resources from two Kubernetes clusters and generate JSON output for comparison.`,
		Run:   runComparison,
	}

	rootCmd.Flags().StringP("output-dir", "o", ".", "Output directory for generated JSON files")
	rootCmd.Flags().BoolP("interactive", "i", true, "Run in interactive mode")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runComparison(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Kubernetes Cluster Resource Comparison Tool")
	fmt.Println("==============================================")
	fmt.Println()

	// Setup and run the comparison
	config, err := setupComparison()
	if err != nil {
		log.Fatalf("Setup failed: %v", err)
	}

	// Fetch resources from both clusters
	if err := fetchResources(config); err != nil {
		log.Fatalf("Failed to fetch resources: %v", err)
	}

	// Generate output files
	if err := generateOutputFiles(config); err != nil {
		log.Fatalf("Failed to generate output files: %v", err)
	}

	fmt.Println("\n🎉 Comparison completed successfully!")
	fmt.Println("📄 Generated files:")
	fmt.Println("   - cluster-a.json")
	fmt.Println("   - cluster-b.json")
	fmt.Println("   - k8s-comparison-report_YYYY-MM-DD_HH-MM-SS.html")
	fmt.Println("💡 Open the HTML report in your browser to view the comparison")
}
