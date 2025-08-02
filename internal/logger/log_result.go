package logger

import (
	"github.com/wcy-dt/ponghub/internal/common"
	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/structures/logger"
	"log"
)

// GetLogs writes check results to JSON file
func GetLogs(checkResult []checker.Checker, maxLogDays int, logPath string) (logger.Logger, error) {
	logResult, err := common.ReadLogs(logPath)
	if err != nil {
		log.Printf("Error loading log data from %s: %v", logPath, err)
		return nil, err
	}

	for _, serviceResult := range checkResult {
		serviceName := serviceResult.Name
		serviceLog, exists := logResult[serviceName]
		if !exists {
			serviceLog = logger.Service{
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
		urlStatusMap, urlTimeMap, urlResponseTimeMap := common.ProcessCheckResult(serviceResult)
		for url, statusList := range urlStatusMap {
			mergedStatus := common.CalcMergedStatus(statusList)
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

	err = common.WriteLogs(logResult, logPath)
	if err != nil {
		log.Printf("Error saving log data to %s: %v", logPath, err)
		return nil, err
	}

	return logResult, nil
}
