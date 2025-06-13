package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

// isGoogleCloudContext checks if a context is a Google Cloud context
func isGoogleCloudContext(contextName string) bool {
	return strings.Contains(contextName, "connectgateway") ||
		strings.Contains(contextName, "gke_") ||
		strings.Contains(contextName, "google")
}

// checkGCloudAuth verifies if gcloud authentication is active
func checkGCloudAuth() error {
	cmd := exec.Command("gcloud", "auth", "list", "--filter=status:ACTIVE", "--format=value(account)")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("gcloud command failed: %w", err)
	}

	if strings.TrimSpace(string(output)) == "" {
		return fmt.Errorf("no active gcloud authentication found")
	}

	return nil
}

// promptGCloudLogin prompts user to authenticate with Google Cloud
func promptGCloudLogin() error {
	fmt.Println("\n⚠️  Google Cloud authentication required!")
	fmt.Println("🔐 Your gcloud credentials have expired or are not set up.")
	fmt.Println("📍 This is required to access Google Kubernetes Engine (GKE) clusters.")

	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Run 'gcloud auth login' now?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("gcloud authentication is required to continue")
	}

	fmt.Println("\n🚀 Opening browser for Google Cloud authentication...")
	fmt.Println("📱 Please complete the authentication process in your browser.")

	cmd := exec.Command("gcloud", "auth", "login")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gcloud auth login failed: %w", err)
	}

	fmt.Println("\n✅ Authentication completed!")

	// Verify authentication worked
	if err := checkGCloudAuth(); err != nil {
		return fmt.Errorf("authentication verification failed: %w", err)
	}

	return nil
}

// ensureGCloudAuth ensures Google Cloud authentication is active for the given context
func ensureGCloudAuth(contextName string) error {
	if !isGoogleCloudContext(contextName) {
		return nil // Not a Google Cloud context, no auth needed
	}

	fmt.Printf("🔍 Checking Google Cloud authentication for context: %s\n", contextName)

	if err := checkGCloudAuth(); err != nil {
		fmt.Printf("⚠️  Authentication issue detected: %v\n", err)
		if err := promptGCloudLogin(); err != nil {
			return err
		}
	} else {
		fmt.Println("✅ Google Cloud authentication is active")
	}

	return nil
}
