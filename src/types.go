package main

// ClusterConfig holds configuration for a single cluster
type ClusterConfig struct {
	Context    string
	Namespaces []string
	Resources  []string
	Data       []map[string]interface{}
}

// ComparisonConfig holds configuration for comparing two clusters
type ComparisonConfig struct {
	ClusterA          ClusterConfig
	ClusterB          ClusterConfig
	OutputDir         string
	ReportTimestamp   string
	CompareNamespaces bool
}
