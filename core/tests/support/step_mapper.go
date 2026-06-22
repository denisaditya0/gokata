package support

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

var serviceCache = make(map[string]reflect.Value)

type Steps struct{}

func RegisterServiceSteps(service string, steps interface{}) {
	serviceCache[service] = reflect.ValueOf(steps)
}

func ValidateServices(services ...string) {
	var missing []string
	for _, s := range services {
		if _, exists := serviceCache[s]; !exists {
			missing = append(missing, s)
		}
	}
	if len(missing) > 0 {
		fmt.Printf("\n❌ Services not registered: %v\n\n", missing)
		fmt.Println("   Add blank imports in 0_suite_test.go:")
		for _, s := range missing {
			fmt.Printf("     _ \"github.com/user/gokata/core/tests/steps/service/%s\"\n", s)
		}
		fmt.Println()
		os.Exit(1)
	}
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
