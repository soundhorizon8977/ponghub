package internal

import (
	"encoding/json"
	"github.com/wcy-dt/ponghub/protos/test_result"
	"log"
	"os"
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

// loadExistingLog loads log data from file or returns empty data
func loadExistingLog(logPath string) (LogData, error) {
	data := make(LogData)

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
func saveLogData(data LogData, logPath string) error {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(logPath, content, 0644)
}

// processCheckResult processes the check results for a service
func processCheckResult(svc CheckResult) (map[string][]test_result.TestResult, map[string]string) {
	urlStatusMap := make(map[string][]test_result.TestResult)
	urlTimeMap := make(map[string]string)

	// Process health checks
	for _, pr := range svc.Health {
		urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online)
		if _, exists := urlTimeMap[pr.URL]; !exists {
			urlTimeMap[pr.URL] = pr.StartTime
		}
	}

	// Process API checks
	for _, pr := range svc.API {
		urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online)
		if _, exists := urlTimeMap[pr.URL]; !exists {
			urlTimeMap[pr.URL] = pr.StartTime
		}
	}

	return urlStatusMap, urlTimeMap
}

// OutputResults writes check results to JSON file
func OutputResults(results []CheckResult, maxLogDays int, logPath string) (LogData, error) {
	logData, err := loadExistingLog(logPath)
	if err != nil {
		log.Printf("Error loading log data from %s: %v", logPath, err)
		return nil, err
	}

	for _, svc := range results {
		serviceName := svc.Name
		serviceLog, exists := logData[serviceName]
		if !exists {
			serviceLog = LogEntry{
				ServiceHistory: HistoryEntryList{},
				PortsData:      make(PortHistory),
			}
		}

		// Update service history
		newHistoryEntry := HistoryEntry{
			Time:   svc.StartTime,
			Status: svc.Online.String(),
		}
		serviceLog.ServiceHistory.AddEntry(newHistoryEntry)
		serviceLog.ServiceHistory.CleanExpiredEntries(maxLogDays)

		// Update port statuses
		urlStatusMap, urlTimeMap := processCheckResult(svc)
		for url, statuses := range urlStatusMap {
			mergedStatus := MergeOnlineStatus(statuses)
			newEntry := HistoryEntry{
				Time:   urlTimeMap[url],
				Status: mergedStatus.String(),
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
