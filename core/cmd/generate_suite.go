// go run cmd/generate_suite.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const module = "github.com/user/gokata/core"

func main() {
	services := findServices("tests/steps/service")

	var imports strings.Builder
	for _, svc := range services {
		imports.WriteString(fmt.Sprintf("\t_ \"%s/tests/steps/service/%s\"\n", module, svc))
	}

	content := fmt.Sprintf(`package features_test

import (
	"testing"

	_ "%s/tests/hooks"

	// Auto-generated: register all services
%s
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFeatures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scenarios Suite")
}
`, module, imports.String())

	os.WriteFile("tests/scenarios/0_suite_test.go", []byte(content), 0644)

	fmt.Printf("Generated 0_suite_test.go with %d services:\n", len(services))
	for _, svc := range services {
		fmt.Printf("  - %s\n", svc)
	}
}

func findServices(dir string) []string {
	var services []string
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() {
			initPath := filepath.Join(dir, entry.Name(), "init.go")
			if _, err := os.Stat(initPath); err == nil {
				services = append(services, entry.Name())
			}
		}
	}
	return services
}
