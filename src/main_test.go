package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestK8sCompare(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Compare Suite")
}
