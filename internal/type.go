package internal

import (
	"github.com/wcy-dt/ponghub/protos/test_result"
	"log"
	"time"
)

// Config defines the configuration structure
type (
	// ServiceConfig defines the configuration for a service, including its health and API ports
	ServiceConfig struct {
		Name    string       `yaml:"name"`
		Health  []PortConfig `yaml:"health"`
		API     []PortConfig `yaml:"api"`
		Timeout int          `yaml:"timeout,omitempty"`
		Retry   int          `yaml:"retry,omitempty"`
	}

	// PortConfig defines the configuration for a port
	PortConfig struct {
		URL           string            `yaml:"url"`
		Method        string            `yaml:"method,omitempty"`
		Headers       map[string]string `yaml:"headers,omitempty"`
		Body          string            `yaml:"body,omitempty"`
		StatusCode    int               `yaml:"status_code,omitempty"`
		ResponseRegex string            `yaml:"response_regex,omitempty"`
	}

	// Config defines the overall configuration structure for the application
	Config struct {
		Services   []ServiceConfig `yaml:"services"`
		Timeout    int             `yaml:"timeout,omitempty"`
		Retry      int             `yaml:"retry,omitempty"`
		MaxLogDays int             `yaml:"max_log_days,omitempty"`
	}
)

// Result defines the structure for the result of checking a service
type (
	// CheckResult defines the structure for the result of checking a service
	CheckResult struct {
		Name          string                 `json:"name"`
		Online        test_result.TestResult `json:"online"`
		Health        []PortResult           `json:"health,omitempty"`
		API           []PortResult           `json:"api,omitempty"`
		StartTime     string                 `json:"start_time"`
		EndTime       string                 `json:"end_time"`
		TotalAttempts int                    `json:"total_attempts"`
		SuccessCount  int                    `json:"success_count"`
	}

	// PortResult defines the structure for the result of checking a port
	PortResult struct {
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

// Data structures for logging and reporting
type (
	// HistoryEntry represents a single history entry
	HistoryEntry struct {
		Time         string `json:"time"`
		Status       string `json:"online"`
		ResponseTime int64  `json:"response_time,omitempty"`
	}

	HistoryEntryList []HistoryEntry
	PortHistory      map[string]HistoryEntryList

	// LogEntry represents log data for a service
	LogEntry struct {
		ServiceHistory HistoryEntryList `json:"service_history"`
		PortsData      PortHistory      `json:"ports"`
	}

	// LogData represents the entire log structure
	LogData map[string]LogEntry

	// ReportEntry represents the result of checking a service
	ReportEntry struct {
		Name         string
		History      HistoryEntryList
		Ports        PortHistory
		Availability float64
	}
)

// CleanExpiredEntries removes entries older than maxDays from the history entry list.
func (hel *HistoryEntryList) CleanExpiredEntries(maxDays int) {
	if maxDays <= 0 {
		log.Println("Max days for cleaning history is not set or invalid, skipping cleaning.")
		return // No cleaning needed if maxDays is not set
	}

	cutoffTime := time.Now().AddDate(0, 0, -maxDays)
	var cleanedHistory HistoryEntryList

	for _, entry := range *hel {
		entryTime, err := time.Parse(time.RFC3339, entry.Time)
		if err != nil {
			log.Printf("Error parsing time %s: %v", entry.Time, err)
			continue // Skip entries with invalid time format
		}
		if entryTime.After(cutoffTime) {
			cleanedHistory = append(cleanedHistory, entry)
		}
	}

	*hel = cleanedHistory
}

// AddEntry adds a new entry to the history entry list.
func (hel *HistoryEntryList) AddEntry(entry HistoryEntry) {
	*hel = append(*hel, entry)
}

// ParseToReportEntries converts LogData to a slice of ReportEntry, calculating availability for each service.
func (ld LogData) ParseToReportEntries() []ReportEntry {
	var reportEntries []ReportEntry
	for svcName, svcData := range ld {
		entryNum := len(svcData.ServiceHistory)
		if entryNum == 0 {
			log.Printf("No history data for service %s", svcName)
			continue // Skip services with no history data
		}

		statusAllEntryNum := 0
		for _, entry := range svcData.ServiceHistory {
			if test_result.IsALL(entry.Status) {
				statusAllEntryNum++
			}
		}
		availability := float64(statusAllEntryNum) / float64(entryNum)

		reportEntries = append(reportEntries, ReportEntry{
			Name:         svcName,
			History:      svcData.ServiceHistory,
			Ports:        svcData.PortsData,
			Availability: availability,
		})
	}
	return reportEntries
}
