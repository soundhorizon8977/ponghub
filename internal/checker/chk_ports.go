package checker

import (
	"errors"
	"fmt"
	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// getHttpMethod converts a string method to an HTTP method constant
func getHttpMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return http.MethodGet
	case "POST":
		return http.MethodPost
	case "PUT":
		return http.MethodPut
	case "DELETE":
		log.Fatalln(errors.New("method not supported"))
	case "HEAD":
		log.Fatalln(errors.New("method not supported"))
	case "PATCH":
		log.Fatalln(errors.New("method not supported"))
	case "OPTIONS":
		log.Fatalln(errors.New("method not supported"))
	case "TRACE":
		log.Fatalln(errors.New("method not supported"))
	case "CONNECT":
		log.Fatalln(errors.New("method not supported"))
	default:
		return http.MethodGet // Default to GET if method is unknown
	}
	return http.MethodGet
}

// isSuccessfulResponse checks if the response from the server is successful based on the configuration
func isSuccessfulResponse(cfg *configure.Port, resp *http.Response, body []byte) bool {
	// responseRegex is set, and the response body does not match the regex
	if cfg.ResponseRegex != "" {
		matched, err := regexp.Match(cfg.ResponseRegex, body)
		if err != nil {
			log.Fatalln("Error parsing regexp:", err)
		}
		if !matched {
			return false
		}
	}

	// statusCode and responseRegex are not set, and the response is OK
	if cfg.StatusCode == 0 && cfg.ResponseRegex == "" && resp.StatusCode == http.StatusOK {
		return true
	}

	// statusCode is not set, and the responseRegex matches
	if cfg.StatusCode == 0 && cfg.ResponseRegex != "" {
		return true
	}

	// statusCode is set, and the response matches the expected status code
	if cfg.StatusCode != 0 && resp.StatusCode == cfg.StatusCode {
		return true
	}

	return false
}

// checkPort checks a single port based on the provided configuration
func checkPort(cfg *configure.Port, timeout int, retryTimes int, svcName string) checker.Port {
	var failures []string
	successCount := 0
	actualAttempts := 0

	var statusCode int
	var responseBody string

	httpMethod := getHttpMethod(cfg.Method)
	responseTime := time.Duration(0)

	// start timer
	start := time.Now()

	for attemptTimes := range retryTimes {
		actualAttempts++
		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
		log.Printf("[%s] %s %s (attempt %d/%d)\n",
			svcName, httpMethod, cfg.URL, attemptTimes+1, retryTimes)

		// build the request
		req, err := http.NewRequest(httpMethod, cfg.URL, nil)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: N/A, Error: %s", err.Error()))
			log.Printf("FAILED - Error: %s", err.Error())
			continue
		}
		for k, v := range cfg.Headers {
			req.Header.Set(k, v)
		}
		if cfg.Body != "" {
			req.Body = io.NopCloser(strings.NewReader(cfg.Body))
		}

		// get the response
		attemptStart := time.Now()
		resp, err := client.Do(req)
		attemptDuration := time.Since(attemptStart)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: N/A, Error: %s", err.Error()))
			log.Printf("FAILED - Error: %s", err.Error())
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			failures = append(failures, fmt.Sprintf("StatusCode: %d, Error: %s", resp.StatusCode, err.Error()))
			log.Printf("FAILED - StatusCode: %d, Error: %s", resp.StatusCode, err.Error())
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body for %s: %v", cfg.URL, err)
			}
			continue
		}
		responseBody = string(body)
		statusCode = resp.StatusCode

		// check the response
		isOnline := isSuccessfulResponse(cfg, resp, body)
		if isOnline {
			successCount++
			if attemptDuration > responseTime {
				responseTime = attemptDuration
			}
			responseBody = ""
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body for %s: %v", cfg.URL, err)
			}
			log.Printf("SUCCESS - %s %s (attempt %d/%d) - Response Time: %d ms, Status Code: %d",
				httpMethod, cfg.URL, attemptTimes+1, retryTimes, attemptDuration.Milliseconds(), resp.StatusCode)
			break
		}
		failures = append(failures, fmt.Sprintf("StatusCode or ResponseRegex mismatch: %d", resp.StatusCode))
		log.Printf("FAILED - StatusCode or ResponseRegex mismatch: %d", resp.StatusCode)
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body for %s: %v", cfg.URL, err)
		}
	}

	// end timer
	end := time.Now()

	return checker.Port{
		URL:           cfg.URL,
		Method:        httpMethod,
		Body:          cfg.Body,
		Online:        getTestResult(successCount, actualAttempts),
		StatusCode:    statusCode,
		StartTime:     start.Format(time.RFC3339),
		EndTime:       end.Format(time.RFC3339),
		ResponseTime:  responseTime,
		TotalAttempts: actualAttempts,
		SuccessCount:  successCount,
		Failures:      failures,
		ResponseBody:  responseBody,
	}
}
