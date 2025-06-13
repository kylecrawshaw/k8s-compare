package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// getKubernetesClient creates a Kubernetes client for the given context
func getKubernetesClient(contextName string) (*kubernetes.Clientset, error) {
	kubeConfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, err := clientcmd.LoadFromFile(kubeConfig)
	if err != nil {
		return nil, err
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{
		CurrentContext: contextName,
	})

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)
}

// getDynamicClient creates a dynamic client and discovery client for the given context
func getDynamicClient(contextName string) (dynamic.Interface, *discovery.DiscoveryClient, error) {
	kubeConfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, err := clientcmd.LoadFromFile(kubeConfig)
	if err != nil {
		return nil, nil, err
	}

	clientConfig := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{
		CurrentContext: contextName,
	})

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, nil, err
	}

	return dynamicClient, discoveryClient, nil
}

// getAvailableContexts returns all available kubectl contexts
func getAvailableContexts() ([]string, error) {
	kubeConfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, err := clientcmd.LoadFromFile(kubeConfig)
	if err != nil {
		return nil, err
	}

	var contexts []string
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}
	sort.Strings(contexts)
	return contexts, nil
}

// getAvailableResourceTypes returns all available resource types for the given context
func getAvailableResourceTypes(contextName string) ([]string, error) {
	client, err := getKubernetesClient(contextName)
	if err != nil {
		return nil, err
	}

	discoveryClient := client.Discovery()
	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	var resources []string
	resourceSet := make(map[string]bool)

	for _, apiResourceList := range apiResourceLists {
		for _, resource := range apiResourceList.APIResources {
			if !strings.Contains(resource.Name, "/") && !resourceSet[resource.Name] {
				resourceSet[resource.Name] = true
				resources = append(resources, resource.Name)
			}
		}
	}

	sort.Strings(resources)
	return resources, nil
}

// fetchClusterResourcesWithContext fetches resources from a cluster with the given context
func fetchClusterResourcesWithContext(contextName string, namespaces []string, resources []string) ([]map[string]interface{}, error) {
	// Add timeout context for operations
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	dynamicClient, discoveryClient, err := getDynamicClient(contextName)
	if err != nil {
		return nil, err
	}

	// Get all API resources
	apiResourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}

	var allItems []runtime.Object

	for _, apiResourceList := range apiResourceLists {
		for _, apiResource := range apiResourceList.APIResources {
			// Check if this resource type is in our selected list
			if !contains(resources, apiResource.Name) {
				continue
			}

			// Skip subresources
			if strings.Contains(apiResource.Name, "/") {
				continue
			}

			// Parse group version
			gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
			if err != nil {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: apiResource.Name,
			}

			// Fetch resources from selected namespaces
			for _, namespace := range namespaces {
				var resourceInterface dynamic.ResourceInterface
				if apiResource.Namespaced {
					resourceInterface = dynamicClient.Resource(gvr).Namespace(namespace)
				} else {
					resourceInterface = dynamicClient.Resource(gvr)
				}

				resources, err := resourceInterface.List(ctx, metav1.ListOptions{})
				if err != nil {
					fmt.Printf("⚠️  Warning: Failed to fetch %s from namespace %s: %v\n", apiResource.Name, namespace, err)
					continue
				}

				for _, item := range resources.Items {
					allItems = append(allItems, &item)
				}
			}
		}
	}

	fmt.Printf("✅ Fetched %d resources from %s\n", len(allItems), contextName)

	var result []map[string]interface{}
	for _, item := range allItems {
		unstructured := item.(*unstructured.Unstructured)
		result = append(result, unstructured.Object)
	}

	return result, nil
}
