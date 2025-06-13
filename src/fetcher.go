package main

import (
	"fmt"
)

// fetchResources fetches resources from both clusters
func fetchResources(config *ComparisonConfig) error {
	fmt.Println("\n📊 Step 4: Fetching resources...")

	var err error

	fmt.Printf("🔍 Fetching resources from Cluster A (%s)...\n", config.ClusterA.Context)
	config.ClusterA.Data, err = fetchClusterResourcesWithContext(config.ClusterA.Context, config.ClusterA.Namespaces, config.ClusterA.Resources)
	if err != nil {
		if isGoogleCloudContext(config.ClusterA.Context) {
			return fmt.Errorf("failed to fetch from Cluster A - this may be due to authentication or network issues with Google Cloud: %w", err)
		}
		return fmt.Errorf("failed to fetch resources from Cluster A: %w", err)
	}
	fmt.Printf("✅ Cluster A: Found %d resources\n", len(config.ClusterA.Data))

	fmt.Printf("🔍 Fetching resources from Cluster B (%s)...\n", config.ClusterB.Context)
	config.ClusterB.Data, err = fetchClusterResourcesWithContext(config.ClusterB.Context, config.ClusterB.Namespaces, config.ClusterB.Resources)
	if err != nil {
		if isGoogleCloudContext(config.ClusterB.Context) {
			return fmt.Errorf("failed to fetch from Cluster B - this may be due to authentication or network issues with Google Cloud: %w", err)
		}
		return fmt.Errorf("failed to fetch resources from Cluster B: %w", err)
	}
	fmt.Printf("✅ Cluster B: Found %d resources\n", len(config.ClusterB.Data))

	return nil
}
