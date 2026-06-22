package support

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"gopkg.in/yaml.v3"
)

// Mode flags
var DryRun = os.Getenv("DRY_RUN") == "true"
var ListMode = os.Getenv("LIST_MODE") == "true"

// --- Test Data Store ---

var testDataStore = make(map[string]map[string]interface{})

func ReadTestDataFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		root := findProjectRoot()
		data, err = os.ReadFile(filepath.Join(root, path))
		if err != nil {
			panic(fmt.Sprintf("failed to load test data from %s: %v", path, err))
		}
	}
	return data
}

func ParseTestData(data []byte) {
	if err := yaml.Unmarshal(data, &testDataStore); err != nil {
		panic(fmt.Sprintf("failed to parse YAML test data: %v", err))
	}
}

func LoadTestDataFromYAML(path string) {
	ParseTestData(ReadTestDataFile(path))
}

func LoadTestDataFromJSON(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load test data from %s: %v", path, err))
	}
	if err := json.Unmarshal(data, &testDataStore); err != nil {
		panic(fmt.Sprintf("failed to parse JSON test data: %v", err))
	}
}

func findProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}

// --- Test Context ---

type TestContext struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	testData map[string]interface{}
}

var ctx = &TestContext{
	data:     make(map[string]interface{}),
	testData: make(map[string]interface{}),
}

func Ctx() *TestContext { return ctx }

// Carryover
func (c *TestContext) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *TestContext) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *TestContext) GetString(key string) string {
	v := c.Get(key)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func (c *TestContext) GetInt(key string) int {
	v := c.Get(key)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
}

func (c *TestContext) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.data[key]
	return exists
}

// Test Data
func (c *TestContext) LoadTestData(scenarioName string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if td, exists := testDataStore[scenarioName]; exists {
		c.testData = td
	} else {
		c.testData = make(map[string]interface{})
	}
}

func (c *TestContext) TD(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.testData[key]
}

func (c *TestContext) TDString(key string) string {
	v := c.TD(key)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func (c *TestContext) TDInt(key string) int {
	v := c.TD(key)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	default:
		return 0
	}
}

func (c *TestContext) TDBool(key string) bool {
	v := c.TD(key)
	if v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}

func (c *TestContext) HasTD(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.testData[key]
	return exists
}

// Clear
func (c *TestContext) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]interface{})
	c.testData = make(map[string]interface{})
}

// --- Fluent Step Builder ---

type Step struct {
	endpoint string
	service  string
}

func Hit(endpoint string) *Step {
	return &Step{endpoint: endpoint}
}

func (s *Step) FromService(service string) *Step {
	s.service = service
	return s
}

func (s *Step) WithCondition(condition string) {
	s.execute(condition)
}

func (s *Step) Execute(condition string) interface{} {
	return s.execute(condition)
}

func (s *Step) execute(condition string) interface{} {
	if DryRun || ListMode {
		fmt.Printf("  - hit %s from %s with condition %s\n", s.endpoint, s.service, condition)
		return nil
	}

	stepName := fmt.Sprintf("hit %s from %s with condition %s", s.endpoint, s.service, condition)
	fmt.Printf("  ○ %s\n", stepName)
	SetCurrentStep(stepName)
	EmitEvent(EventStep, map[string]string{"name": stepName})

	svc, exists := serviceCache[s.service]
	if !exists {
		panic(fmt.Sprintf("Service '%s' not registered", s.service))
	}

	funcName := "Execute" + toPascalCase(s.endpoint)
	method := svc.MethodByName(funcName)
	if !method.IsValid() {
		panic(fmt.Sprintf("Method '%s' not found in service '%s'", funcName, s.service))
	}

	results := method.Call([]reflect.Value{reflect.ValueOf(condition)})
	if len(results) > 0 {
		return results[0].Interface()
	}
	return nil
}
