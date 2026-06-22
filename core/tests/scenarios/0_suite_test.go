package features_test

import (
	"testing"

	_ "github.com/user/gokata/core/tests/hooks"

	// Auto-generated: register all services
	_ "github.com/user/gokata/core/tests/steps/service/auth"
	_ "github.com/user/gokata/core/tests/steps/service/products"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFeatures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scenarios Suite")
}
