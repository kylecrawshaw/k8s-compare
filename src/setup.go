package main

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// setupComparison handles the interactive setup process
func setupComparison() (*ComparisonConfig, error) {
	// Get available contexts
	contexts, err := getAvailableContexts()
	if err != nil {
		return nil, fmt.Errorf("failed to get contexts: %w", err)
	}

	if len(contexts) < 2 {
		return nil, fmt.Errorf("need at least 2 contexts, found %d", len(contexts))
	}

	config := &ComparisonConfig{}

	// Select contexts
	fmt.Println("📍 Step 1: Select Kubernetes contexts")
	config.ClusterA.Context, err = selectFromList("Select Cluster A context:", contexts)
	if err != nil {
		return nil, err
	}

	remainingContexts := removeFromSlice(contexts, config.ClusterA.Context)
	config.ClusterB.Context, err = selectFromList("Select Cluster B context:", remainingContexts)
	if err != nil {
		return nil, err
	}

	// Early authentication check for Google Cloud contexts
	fmt.Println("\n🔐 Checking authentication for selected contexts...")

	if err := ensureGCloudAuth(config.ClusterA.Context); err != nil {
		return nil, fmt.Errorf("authentication failed for Cluster A (%s): %w", config.ClusterA.Context, err)
	}

	if err := ensureGCloudAuth(config.ClusterB.Context); err != nil {
		return nil, fmt.Errorf("authentication failed for Cluster B (%s): %w", config.ClusterB.Context, err)
	}

	// Select namespaces for each cluster
	fmt.Println("\n🏠 Step 2: Select namespaces")
	config.ClusterA.Namespaces, err = selectNamespaces(config.ClusterA.Context, "Cluster A")
	if err != nil {
		return nil, err
	}

	config.ClusterB.Namespaces, err = selectNamespaces(config.ClusterB.Context, "Cluster B")
	if err != nil {
		return nil, err
	}

	// Select resource types
	fmt.Println("\n📦 Step 3: Select resource types")
	availableResources, err := getAvailableResourceTypes(config.ClusterA.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to get available resource types: %w", err)
	}

	config.ClusterA.Resources, err = selectMultipleFromList("Select resource types to compare:", availableResources)
	if err != nil {
		return nil, err
	}
	config.ClusterB.Resources = config.ClusterA.Resources

	return config, nil
}

// selectFromList presents a single-select list to the user
func selectFromList(title string, items []string) (string, error) {
	var selected string

	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item, item)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(title).
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// selectMultipleFromList presents a multi-select list to the user
func selectMultipleFromList(title string, items []string) ([]string, error) {
	var selected []string

	// Reorder items to show common resources first
	reorderedItems := reorderResourcesByPriority(items)

	options := make([]huh.Option[string], len(reorderedItems))
	for i, item := range reorderedItems {
		options[i] = huh.NewOption(item, item)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(title).
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	return selected, err
}

// reorderResourcesByPriority puts common resources first in the list
func reorderResourcesByPriority(items []string) []string {
	// Common resource types that should appear first
	commonResources := []string{
		"pods", "services", "deployments", "statefulsets", "configmaps",
		"secrets", "persistentvolumeclaims", "jobs", "cronjobs", "replicasets",
		"daemonsets", "ingresses", "networkpolicies", "persistentvolumes",
	}

	// Filter to only include items that exist in the available list
	var filteredCommon []string
	for _, common := range commonResources {
		for _, item := range items {
			if strings.Contains(item, common) {
				filteredCommon = append(filteredCommon, item)
				break
			}
		}
	}

	// Prepend common resources, then add remaining
	var reorderedItems []string
	reorderedItems = append(reorderedItems, filteredCommon...)

	for _, item := range items {
		found := false
		for _, common := range filteredCommon {
			if item == common {
				found = true
				break
			}
		}
		if !found {
			reorderedItems = append(reorderedItems, item)
		}
	}

	return reorderedItems
}

// selectNamespaces handles namespace selection for a cluster
func selectNamespaces(contextName, clusterName string) ([]string, error) {
	// Ensure Google Cloud authentication if needed
	if err := ensureGCloudAuth(contextName); err != nil {
		return nil, fmt.Errorf("google cloud authentication failed: %w", err)
	}

	client, err := getKubernetesClient(contextName)
	if err != nil {
		if isGoogleCloudContext(contextName) && (strings.Contains(err.Error(), "gke-gcloud-auth-plugin") ||
			strings.Contains(err.Error(), "credential") ||
			strings.Contains(err.Error(), "auth")) {
			fmt.Println("\n🔄 Authentication issue detected, attempting to refresh credentials...")
			if authErr := promptGCloudLogin(); authErr != nil {
				return nil, fmt.Errorf("failed to authenticate with Google Cloud: %w", authErr)
			}
			// Retry after authentication
			client, err = getKubernetesClient(contextName)
			if err != nil {
				return nil, fmt.Errorf("failed to connect even after authentication: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to create Kubernetes client for context %s: %w", contextName, err)
		}
	}

	fmt.Printf("📋 Fetching namespaces from %s...\n", clusterName)
	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if isGoogleCloudContext(contextName) && (strings.Contains(err.Error(), "credential") ||
			strings.Contains(err.Error(), "auth") ||
			strings.Contains(err.Error(), "token")) {
			fmt.Println("\n🔄 Token expired, refreshing authentication...")
			if authErr := promptGCloudLogin(); authErr != nil {
				return nil, fmt.Errorf("failed to refresh authentication: %w", authErr)
			}
			// Retry after re-authentication
			client, err = getKubernetesClient(contextName)
			if err != nil {
				return nil, fmt.Errorf("failed to reconnect after auth refresh: %w", err)
			}
			namespaces, err = client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to list namespaces even after re-authentication: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}
	}

	var nsNames []string
	for _, ns := range namespaces.Items {
		nsNames = append(nsNames, ns.Name)
	}
	sort.Strings(nsNames)

	selectedNs, err := selectMultipleFromList(fmt.Sprintf("Select namespaces for %s:", clusterName), nsNames)
	if err != nil {
		return nil, err
	}

	if len(selectedNs) == 0 {
		return nil, fmt.Errorf("no namespaces selected")
	}

	return selectedNs, nil
}
