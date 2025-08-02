package checker

import (
	"github.com/wcy-dt/ponghub/internal/types/types/test_result"
	"time"
)

// Result defines the structure for the result of checking a service
type (
	// Checker defines the structure for the result of checking a service
	Checker struct {
		Name          string                 `json:"name"`
		Online        test_result.TestResult `json:"online"`
		Health        []Port                 `json:"health,omitempty"`
		API           []Port                 `json:"api,omitempty"`
		StartTime     string                 `json:"start_time"`
		EndTime       string                 `json:"end_time"`
		TotalAttempts int                    `json:"total_attempts"`
		SuccessCount  int                    `json:"success_count"`
	}

	// Port defines the structure for the result of checking a port
	Port struct {
		URL           string                 `json:"url"`
		Method        string                 `json:"method"`
		Body          string                 `json:"body,omitempty"`
		Online        test_result.TestResult `json:"online"`
		StatusCode    int                    `json:"status_code,omitempty"`
		StartTime     string                 `json:"start_time"`
		EndTime       string                 `json:"end_time"`
		ResponseTime  time.Duration          `json:"response_time"`
		TotalAttempts int                    `json:"total_attempts"`
		SuccessCount  int                    `json:"success_count"`
		Failures      []string               `json:"failures,omitempty"`
		ResponseBody  string                 `json:"response_body,omitempty"`
	}
)
