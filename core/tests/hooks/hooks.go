package hooks

import (
	"fmt"
	"os"

	"github.com/user/gokata/core/tests/support"
	. "github.com/onsi/ginkgo/v2"
)

var _ = SynchronizedBeforeSuite(
	func() []byte {
		fmt.Println("🚀 Starting Test Suite")
		support.ValidateServices("auth", "products")
		support.InitEvents()

		mode := os.Getenv("DATA_MODE")
		name := os.Getenv("DATA_NAME")
		env := os.Getenv("DATA_ENV")

		if mode == "" || name == "" || env == "" {
			fmt.Println("⚠️  No test data config (DATA_MODE/DATA_NAME/DATA_ENV), running without test data")
			return nil
		}

		path := fmt.Sprintf("tests/data/%s/%s/%s.yaml", mode, name, env)
		data := support.ReadTestDataFile(path)
		fmt.Printf("📂 Loaded test data: %s\n", path)

		support.EmitEvent(support.EventSuiteStart, map[string]interface{}{
			"mode": mode, "name": name, "env": env,
		})

		return data
	},
	func(data []byte) {
		if len(data) > 0 {
			support.ParseTestData(data)
		}
	},
)

var _ = AfterSuite(func() {
	fmt.Println("✅ Finished Test Suite")
	support.EmitEvent(support.EventSuiteEnd, nil)
	support.CloseEvents()
	support.CloseHTTPConnections()
})
