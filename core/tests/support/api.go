package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/jmespath/go-jmespath"
	. "github.com/onsi/gomega"
)

// --- HTTP Client ---

var (
	httpClient     *http.Client
	httpClientOnce sync.Once
)

func getHTTPClient() *http.Client {
	httpClientOnce.Do(func() {
		httpClient = &http.Client{
			Timeout: GetTimeout(),
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
	})
	return httpClient
}

// CloseHTTPConnections closes all idle HTTP connections
func CloseHTTPConnections() {
	if httpClient != nil {
		httpClient.CloseIdleConnections()
	}
}

// --- HTTP Response ---

type HTTPResponse struct {
	StatusCode int
	Body       map[string]interface{}
	RawBody    []byte
	Error      error
}

// --- Request Builder ---

func BuildRequest(method, endpoint string, body interface{}) *http.Request {
	url := GetConfig().BaseURL + endpoint

	var reqBody *bytes.Buffer
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// --- HTTP Executor (with retry + logging) ---

var HTTPLog = true
var currentStep string

func SetCurrentStep(step string) { currentStep = step }

func ExecuteHTTPRequest(req *http.Request) *HTTPResponse {
	client := getHTTPClient()
	maxRetry := GetConfig().Retry

	// Read request body for logging/retry
	var reqBodyBytes []byte
	var reqBody interface{}
	if req.Body != nil {
		reqBodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		if len(reqBodyBytes) > 0 {
			json.Unmarshal(reqBodyBytes, &reqBody)
		}
	}

	if HTTPLog {
		fmt.Printf("    ▶ %s %s\n", req.Method, req.URL.String())
		if len(reqBodyBytes) > 0 {
			fmt.Printf("      Request: %s\n", string(reqBodyBytes))
		}
	}

	var resp *http.Response
	var err error
	start := time.Now()

	// Retry loop
	for attempt := 0; attempt <= maxRetry; attempt++ {
		if attempt > 0 {
			// Reset body for retry
			req.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			if HTTPLog {
				fmt.Printf("      ↻ Retry %d/%d\n", attempt, maxRetry)
			}
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		resp, err = client.Do(req)
		if err == nil {
			break
		}
	}

	if err != nil {
		if HTTPLog {
			fmt.Printf("      ✗ Error: %v\n", err)
		}
		return &HTTPResponse{Error: err}
	}
	defer resp.Body.Close()

	rawBody, _ := io.ReadAll(resp.Body)
	duration := time.Since(start)

	var body map[string]interface{}
	json.Unmarshal(rawBody, &body)

	if HTTPLog {
		fmt.Printf("    ◀ %d (%s)\n", resp.StatusCode, duration.Round(time.Millisecond))
		if len(rawBody) > 0 && len(rawBody) < 500 {
			fmt.Printf("      Response: %s\n", string(rawBody))
		} else if len(rawBody) >= 500 {
			fmt.Printf("      Response: %s...\n", string(rawBody[:200]))
		}
	}

	// Structured log for report
	logEntry := StepLog{
		Step:     currentStep,
		Method:   req.Method,
		URL:      req.URL.String(),
		Request:  reqBody,
		Status:   resp.StatusCode,
		Response: body,
		Duration: duration.Round(time.Millisecond).String(),
	}
	scenarioLogs = append(scenarioLogs, logEntry)
	EmitEvent(EventHTTP, logEntry)

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		RawBody:    rawBody,
	}
}

// --- Structured Logs ---

type StepLog struct {
	Step     string      `json:"step"`
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Request  interface{} `json:"request,omitempty"`
	Status   int         `json:"status"`
	Response interface{} `json:"response,omitempty"`
	Duration string      `json:"duration"`
}

var scenarioLogs []StepLog

func GetScenarioLogs() []StepLog {
	logs := scenarioLogs
	scenarioLogs = nil
	return logs
}

func ResetScenarioLogs() { scenarioLogs = nil }

// --- Validators (Gomega) ---

func ValidateStatusCode(response *HTTPResponse, expectedStatus int) {
	Expect(response.StatusCode).To(Equal(expectedStatus),
		fmt.Sprintf("Expected status %d, got %d", expectedStatus, response.StatusCode))
}

func ValidateResponseBody(response *HTTPResponse, field string) {
	Expect(response.Body).To(HaveKey(field),
		fmt.Sprintf("Response missing field: %s", field))
}

func ValidateErrorMessage(response *HTTPResponse, expectedMessage string) {
	Expect(response.Error).NotTo(BeNil(), "Expected error but got none")
	Expect(response.Error.Error()).To(Equal(expectedMessage))
}

func ValidateBodyField(response *HTTPResponse, field string, expectedValue interface{}) {
	Expect(response.Body).To(HaveKey(field))
	Expect(response.Body[field]).To(Equal(expectedValue))
}

// --- JMESPath ---

func JMESPathQuery(data interface{}, query string) (interface{}, error) {
	return jmespath.Search(query, data)
}

func JMESPathQueryString(data interface{}, query string) (string, error) {
	result, err := jmespath.Search(query, data)
	if err != nil {
		return "", err
	}
	switch v := result.(type) {
	case string:
		return v, nil
	case nil:
		return "", fmt.Errorf("query returned nil")
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func JMESPathQueryInt(data interface{}, query string) (int, error) {
	result, err := jmespath.Search(query, data)
	if err != nil {
		return 0, err
	}
	switch v := result.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", result)
	}
}

// --- JMESPath Validators (Gomega) ---

func ValidateJMESPathValue(data interface{}, query string, expectedValue interface{}) {
	result, err := jmespath.Search(query, data)
	Expect(err).NotTo(HaveOccurred(), "JMESPath query failed")
	Expect(result).To(Equal(expectedValue), fmt.Sprintf("JMESPath '%s'", query))
}

func ValidateJMESPathExists(data interface{}, query string) {
	result, err := jmespath.Search(query, data)
	Expect(err).NotTo(HaveOccurred())
	Expect(result).NotTo(BeNil(), fmt.Sprintf("JMESPath '%s' returned nil", query))
}
