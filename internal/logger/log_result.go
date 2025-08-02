package logger

import (
	"encoding/json"
	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/structures/logger"
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
	"log"
	"os"
	"time"
)

// GetLogs writes check results to JSON file
func GetLogs(checkResult []checker.Checker, maxLogDays int, logPath string) (logger.Logger, error) {
	logResult, err := ReadLogs(logPath)
	if err != nil {
		log.Printf("Error loading log data from %s: %v", logPath, err)
		return nil, err
	}

	for _, serviceResult := range checkResult {
		serviceName := serviceResult.Name
		serviceLog, exists := logResult[serviceName]
		if !exists {
			serviceLog = logger.Entry{
				ServiceHistory: logger.History{},
				Endpoints:      make(logger.Endpoints),
			}
		}

		// Update service history
		newServiceHistoryEntry := logger.HistoryEntry{
			Time:   serviceResult.StartTime, // Use StartTime for the history entry
			Status: serviceResult.Status.String(),
		}
		serviceLog.ServiceHistory = serviceLog.ServiceHistory.AddEntry(newServiceHistoryEntry)
		serviceLog.ServiceHistory = serviceLog.ServiceHistory.CleanExpiredEntries(maxLogDays)

		// Update port statusList
		urlStatusMap, urlTimeMap, urlResponseTimeMap := processCheckResult(serviceResult)
		for url, statusList := range urlStatusMap {
			mergedStatus := calcMergedStatus(statusList)
			newEndpointHistoryEntry := logger.HistoryEntry{
				Time:         urlTimeMap[url],
				Status:       mergedStatus.String(),
				ResponseTime: int(urlResponseTimeMap[url].Milliseconds()),
			}

			tmp := serviceLog.Endpoints[url]
			tmp = tmp.AddEntry(newEndpointHistoryEntry)
			tmp = tmp.CleanExpiredEntries(maxLogDays)
			serviceLog.Endpoints[url] = tmp
		}

		logResult[serviceName] = serviceLog
	}

	err = WriteLogs(logResult, logPath)
	if err != nil {
		log.Printf("Error saving log data to %s: %v", logPath, err)
		return nil, err
	}

	return logResult, nil
}

// calcMergedStatus merges multiple statuses into a single status
func calcMergedStatus(statusList []chk_result.CheckResult) chk_result.CheckResult {
	if len(statusList) == 0 {
		return chk_result.NONE
	}

	hasNone, hasAll := false, false
	for _, s := range statusList {
		switch s {
		case chk_result.NONE:
			hasNone = true
		case chk_result.ALL:
			hasAll = true
		}
	}

	switch {
	case hasNone && !hasAll:
		return chk_result.NONE
	case !hasNone && hasAll:
		return chk_result.ALL
	default:
		return chk_result.PART
	}
}

// ReadLogs loads log data from file or returns empty data
func ReadLogs(logPath string) (logger.Logger, error) {
	logResult := make(logger.Logger)

	logContent, err := os.ReadFile(logPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return logResult, nil
	}

	if err := json.Unmarshal(logContent, &logResult); err != nil {
		return nil, err
	}
	return logResult, nil
}

// WriteLogs writes log data to file
func WriteLogs(logResult logger.Logger, logPath string) error {
	logContent, err := json.MarshalIndent(logResult, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(logPath, logContent, 0644)
}

// processCheckResult processes the check results for a service
func processCheckResult(serviceResult checker.Checker) (map[string][]chk_result.CheckResult, map[string]string, map[string]time.Duration) {
	urlStatusMap := make(map[string][]chk_result.CheckResult)
	urlTimeMap := make(map[string]string)
	urlResponseTimeMap := make(map[string]time.Duration)

	// Process Endpoints checks
	for _, endpoint := range serviceResult.Endpoints {
		urlStatusMap[endpoint.URL] = append(urlStatusMap[endpoint.URL], endpoint.Status)

		if _, exists := urlTimeMap[endpoint.URL]; !exists {
			urlTimeMap[endpoint.URL] = endpoint.StartTime
		} else if endpoint.StartTime < urlTimeMap[endpoint.URL] {
			urlTimeMap[endpoint.URL] = endpoint.StartTime
		}

		if _, exists := urlResponseTimeMap[endpoint.URL]; !exists {
			urlResponseTimeMap[endpoint.URL] = endpoint.ResponseTime
		} else if endpoint.ResponseTime > urlResponseTimeMap[endpoint.URL] {
			urlResponseTimeMap[endpoint.URL] = endpoint.ResponseTime
		}
	}

	return urlStatusMap, urlTimeMap, urlResponseTimeMap
}
