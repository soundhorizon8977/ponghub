package checker

import (
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
	"time"
)

// Result defines the structure for the result of checking a service
type (
	// Checker defines the structure for the result of checking a service
	Checker struct {
		Name       string                 `json:"name"`
		Status     chk_result.CheckResult `json:"status"`
		Endpoints  []Endpoint             `json:"endpoints,omitempty"`
		StartTime  string                 `json:"start_time"`
		EndTime    string                 `json:"end_time"`
		AttemptNum int                    `json:"attempt_num"`
		SuccessNum int                    `json:"success_num"`
	}

	// Endpoint defines the structure for the result of checking a port
	Endpoint struct {
		URL            string                 `json:"url"`
		Method         string                 `json:"method"`
		Body           string                 `json:"body,omitempty"`
		Status         chk_result.CheckResult `json:"status"`
		StatusCode     int                    `json:"status_code,omitempty"`
		StartTime      string                 `json:"start_time"`
		EndTime        string                 `json:"end_time"`
		ResponseTime   time.Duration          `json:"response_time"`
		AttemptNum     int                    `json:"attempt_num"`
		SuccessNum     int                    `json:"success_num"`
		FailureDetails []string               `json:"failure_details,omitempty"`
		ResponseBody   string                 `json:"response_body,omitempty"`
	}
)
