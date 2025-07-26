package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wcy-dt/ponghub/protos/default_config"
	"github.com/wcy-dt/ponghub/protos/test_result"
	"html/template"
	"os"
)

type (
	ServiceHistory struct {
		Status string
		Time   string
	}

	PortHistory struct {
		URL    string
		Time   string
		Status string
	}

	ServiceResult struct {
		Name         string
		History      []ServiceHistory
		Ports        map[string][]PortHistory
		Availability float64
	}

	LogData map[string]ServiceData

	ServiceData struct {
		ServiceHistory []map[string]any `json:"service_history"`
		Ports          PortCollection   `json:"ports"`
	}

	PortCollection map[string][]map[string]any
)

// parseLogData parse the log data into a structured format
func parseLogData(logData LogData) ([]ServiceResult, string, error) {
	if len(logData) == 0 {
		return nil, "", errors.New("log data is empty")
	}

	var (
		results    []ServiceResult
		latestTime string
	)

	for svcName, svcData := range logData {
		history, svcLatestTime, availability := parseServiceHistory(svcData.ServiceHistory)
		ports, portsLatestTime := parsePortData(svcData.Ports)

		// determine the latest time for this service
		serviceLatest := latestOf(svcLatestTime, portsLatestTime)
		if serviceLatest > latestTime {
			latestTime = serviceLatest
		}

		results = append(results, ServiceResult{
			Name:         svcName,
			History:      history,
			Ports:        ports,
			Availability: availability,
		})
	}

	return results, latestTime, nil
}

// parseServiceHistory parses the service history data
func parseServiceHistory(items []map[string]any) ([]ServiceHistory, string, float64) {
	var (
		histories    []ServiceHistory
		maxTime      string
		successCount int
	)

	for _, item := range items {
		status, ok1 := item["online"].(string)
		timeStr, ok2 := item["time"].(string)

		if !ok1 || !ok2 {
			continue
		}

		history := ServiceHistory{
			Status: status,
			Time:   timeStr,
		}
		histories = append(histories, history)

		if test_result.IsALL(status) {
			successCount++
		}

		if timeStr > maxTime {
			maxTime = timeStr
		}
	}

	availability := 0.0
	if count := len(histories); count > 0 {
		availability = float64(successCount) / float64(count)
	}

	return histories, maxTime, availability
}

// parsePortData parses the port data from the log
func parsePortData(ports PortCollection) (map[string][]PortHistory, string) {
	result := make(map[string][]PortHistory)
	var maxTime string

	for portURL, historyItems := range ports {
		var entries []PortHistory

		for _, item := range historyItems {
			status, ok1 := item["online"].(string)
			timeStr, ok2 := item["time"].(string)

			if !ok1 || !ok2 {
				continue
			}

			entry := PortHistory{
				URL:    portURL,
				Status: status,
				Time:   timeStr,
			}
			entries = append(entries, entry)

			if timeStr > maxTime {
				maxTime = timeStr
			}
		}

		if len(entries) > 0 {
			result[portURL] = entries
		}
	}

	return result, maxTime
}

// latestOf returns the latest time between two strings
func latestOf(a, b string) string {
	if a > b {
		return a
	}
	return b
}

// GenerateReport generates an HTML report from the log data
func GenerateReport(logPath, outPath string) error {
	data, err := os.ReadFile(logPath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	var logData LogData
	if err := json.Unmarshal(data, &logData); err != nil {
		return fmt.Errorf("failed to parse log data: %w", err)
	}

	results, latestTime, err := parseLogData(logData)
	if err != nil {
		return fmt.Errorf("processing log data: %w", err)
	}

	tmpl, err := template.New("report.html").
		Funcs(createTemplateFunc()).
		ParseFiles(default_config.GetTemplatePath())
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	outputFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			fmt.Printf("Error closing output file: %v\n", err)
		}
	}(outputFile)

	if err := tmpl.Execute(outputFile, map[string]interface{}{
		"Results":    results,
		"UpdateTime": latestTime,
	}); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// createTemplateFunc defines custom template functions for the report
func createTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"until": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
		},
	}
}
