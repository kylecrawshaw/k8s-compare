package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "k8s-compare",
		Short: "Compare resources between two Kubernetes clusters",
		Long:  `A CLI tool to interactively select contexts, namespaces, and resources from two Kubernetes clusters and generate JSON output for comparison.`,
		Run:   runComparison,
	}

	rootCmd.Flags().StringP("output-dir", "o", "reports", "Output directory for generated JSON files")
	rootCmd.Flags().BoolP("interactive", "i", true, "Run in interactive mode")
	rootCmd.Flags().BoolP("compare-namespaces", "c", true, "Compare namespaces")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runComparison(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Kubernetes Cluster Resource Comparison Tool")
	fmt.Println("==============================================")
	fmt.Println()

	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		log.Fatalf("Failed to get output directory: %v", err)
	}

	compareNamespaces, err := cmd.Flags().GetBool("compare-namespaces")
	if err != nil {
		log.Fatalf("Failed to get compare namespaces flag: %v", err)
	}

	// Setup and run the comparison
	config, err := setupComparison(outputDir, compareNamespaces)
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
	fmt.Printf("   - %s/cluster-a-%s.json\n", config.OutputDir, config.ReportTimestamp)
	fmt.Printf("   - %s/cluster-b-%s.json\n", config.OutputDir, config.ReportTimestamp)
	fmt.Printf("   - %s/k8s-comparison-report_%s.html\n", config.OutputDir, config.ReportTimestamp)
	fmt.Println("💡 Open the HTML report in your browser to view the comparison")
	fmt.Printf("   👉 Example: open %s/k8s-comparison-report_%s.html\n", config.OutputDir, config.ReportTimestamp)

	interactive, _ := cmd.Flags().GetBool("interactive")
	if interactive {
		reportFile := fmt.Sprintf("%s/k8s-comparison-report_%s.html", config.OutputDir, config.ReportTimestamp)
		var openNow bool
		huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Title("Would you like to open the HTML report now?").Value(&openNow),
			),
		).Run()
		if openNow {
			var cmdOpen *exec.Cmd
			switch runtime.GOOS {
			case "darwin":
				cmdOpen = exec.Command("open", reportFile)
			case "linux":
				cmdOpen = exec.Command("xdg-open", reportFile)
			case "windows":
				cmdOpen = exec.Command("cmd", "/c", "start", reportFile)
			default:
				fmt.Println("Cannot determine how to open files on this OS.")
				return
			}
			err := cmdOpen.Start()
			if err != nil {
				fmt.Printf("Failed to open report: %v\n", err)
			}
		}
	}
}

func checkedAttr(val bool) string {
	if val {
		return " checked"
	}
	return ""
}
