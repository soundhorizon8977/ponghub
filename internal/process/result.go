package process

import (
	"encoding/json"
	"github.com/wcy-dt/ponghub/internal/types"
	"github.com/wcy-dt/ponghub/internal/types/test_result"
	"log"
	"os"
	"time"
)

// MergeOnlineStatus merges multiple statuses into a single status
func MergeOnlineStatus(statuses []test_result.TestResult) test_result.TestResult {
	if len(statuses) == 0 {
		return test_result.NONE
	}

	hasNone, hasAll := false, false
	for _, s := range statuses {
		switch s {
		case test_result.NONE:
			hasNone = true
		case test_result.ALL:
			hasAll = true
		}
	}

	switch {
	case hasNone && !hasAll:
		return test_result.NONE
	case !hasNone && hasAll:
		return test_result.ALL
	default:
		return test_result.PART
	}
}

// LoadExistingLog loads log data from file or returns empty data
func LoadExistingLog(logPath string) (types.LogData, error) {
	data := make(types.LogData)

	content, err := os.ReadFile(logPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return data, nil
	}

	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// saveLogData writes log data to file
func saveLogData(data types.LogData, logPath string) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(logPath, content, 0644)
}

// processCheckResult processes the check results for a service
func processCheckResult(svc types.CheckResult) (map[string][]test_result.TestResult, map[string]string, map[string]time.Duration) {
	urlStatusMap := make(map[string][]test_result.TestResult)
	urlTimeMap := make(map[string]string)
	urlResponseTimeMap := make(map[string]time.Duration)

	// Process health checks
	for _, pr := range svc.Health {
		urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online)
		if _, exists := urlTimeMap[pr.URL]; !exists {
			urlTimeMap[pr.URL] = pr.StartTime
		}
		if _, exists := urlResponseTimeMap[pr.URL]; !exists {
			urlResponseTimeMap[pr.URL] = pr.ResponseTime
		} else if pr.ResponseTime > urlResponseTimeMap[pr.URL] {
			urlResponseTimeMap[pr.URL] = pr.ResponseTime
		}
	}

	// Process API checks
	for _, pr := range svc.API {
		urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online)
		if _, exists := urlTimeMap[pr.URL]; !exists {
			urlTimeMap[pr.URL] = pr.StartTime
		}
		if _, exists := urlResponseTimeMap[pr.URL]; !exists {
			urlResponseTimeMap[pr.URL] = pr.ResponseTime
		} else if pr.ResponseTime > urlResponseTimeMap[pr.URL] {
			urlResponseTimeMap[pr.URL] = pr.ResponseTime
		}
	}

	return urlStatusMap, urlTimeMap, urlResponseTimeMap
}

// OutputResults writes check results to JSON file
func OutputResults(results []types.CheckResult, maxLogDays int, logPath string) (types.LogData, error) {
	logData, err := LoadExistingLog(logPath)
	if err != nil {
		log.Printf("Error loading log data from %s: %v", logPath, err)
		return nil, err
	}

	for _, svc := range results {
		serviceName := svc.Name
		serviceLog, exists := logData[serviceName]
		if !exists {
			serviceLog = types.LogEntry{
				ServiceHistory: types.HistoryEntryList{},
				PortsData:      make(types.PortHistory),
			}
		}

		// Update service history
		newHistoryEntry := types.HistoryEntry{
			Time:   svc.StartTime,
			Status: svc.Online.String(),
		}
		serviceLog.ServiceHistory.AddEntry(newHistoryEntry)
		serviceLog.ServiceHistory.CleanExpiredEntries(maxLogDays)

		// Update port statuses
		urlStatusMap, urlTimeMap, urlResponseTimeMap := processCheckResult(svc)
		for url, statuses := range urlStatusMap {
			mergedStatus := MergeOnlineStatus(statuses)
			newEntry := types.HistoryEntry{
				Time:         urlTimeMap[url],
				Status:       mergedStatus.String(),
				ResponseTime: int(urlResponseTimeMap[url].Milliseconds()),
			}

			tmp := serviceLog.PortsData[url]
			tmp.AddEntry(newEntry)
			tmp.CleanExpiredEntries(maxLogDays)
			serviceLog.PortsData[url] = tmp
		}

		logData[serviceName] = serviceLog
	}

	err = saveLogData(logData, logPath)
	if err != nil {
		log.Printf("Error saving log data to %s: %v", logPath, err)
		return nil, err
	}

	return logData, nil
}
