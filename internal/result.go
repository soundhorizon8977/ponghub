package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wcy-dt/ponghub/protos/test_result"
)

// MergeOnlineStatus merges a list of online statuses into a single status
func MergeOnlineStatus(statusList []test_result.TestResult) test_result.TestResult {
	if len(statusList) == 0 {
		return test_result.NONE
	}

	hasNone, hasAll := false, false
	for _, s := range statusList {
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

// OutputResults writes the check results to a JSON file and updates the log file
func OutputResults(results []CheckResult, maxLogDays int, logPath string) error {
	// Get existing log data or create a new map
	var logData = make(map[string]map[string]any)
	if b, err := os.ReadFile(logPath); err == nil {
		if err := json.Unmarshal(b, &logData); err != nil {
			log.Fatalln("Failed to read existing log file:", err)
		}
	}

	// last update time
	now := time.Now()

	for _, svc := range results {
		// check if service already exists in logData, otherwise initialize it
		if _, ok := logData[svc.Name]; !ok {
			logData[svc.Name] = map[string]any{
				"service_history": []any{},
				"ports":           map[string][]any{},
			}
		}

		// Handle service_history type
		svcHistoryRaw := logData[svc.Name]["service_history"]
		var svcHistory []map[string]string
		switch v := svcHistoryRaw.(type) {
		case []any:
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					entry := map[string]string{}
					for k, val := range m {
						entry[k] = fmt.Sprintf("%v", val)
					}
					svcHistory = append(svcHistory, entry)
				}
			}
		case []map[string]string:
			svcHistory = v
		}
		svcHistory = append(svcHistory, map[string]string{
			"time":   svc.StartTime,
			"online": svc.Online.String(),
		})
		// Clean up timeout records
		var filteredSvcHistory []map[string]string
		for _, entry := range svcHistory {
			t, err := time.Parse(time.RFC3339, entry["time"])
			if err == nil && now.Sub(t).Hours() <= float64(maxLogDays*24) {
				filteredSvcHistory = append(filteredSvcHistory, entry)
			}
		}
		logData[svc.Name]["service_history"] = filteredSvcHistory

		// Handle ports type
		portsRaw := logData[svc.Name]["ports"]
		portsMap := map[string][]map[string]string{}
		switch v := portsRaw.(type) {
		case map[string]any:
			for url, arr := range v {
				var portHistory []map[string]string
				if arrList, ok := arr.([]any); ok {
					for _, item := range arrList {
						if m, ok := item.(map[string]any); ok {
							entry := map[string]string{}
							for k, val := range m {
								entry[k] = fmt.Sprintf("%v", val)
							}
							portHistory = append(portHistory, entry)
						}
					}
				}
				portsMap[url] = portHistory
			}
		case map[string][]map[string]string:
			portsMap = v
		}
		// Only record one port entry for each unique URL per complete run
		urlStatusMap := map[string][]string{}
		urlTimeMap := map[string]string{}
		for _, pr := range svc.Health {
			urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online.String())
			if urlTimeMap[pr.URL] == "" {
				urlTimeMap[pr.URL] = pr.StartTime
			}
		}
		for _, pr := range svc.API {
			urlStatusMap[pr.URL] = append(urlStatusMap[pr.URL], pr.Online.String())
			if urlTimeMap[pr.URL] == "" {
				urlTimeMap[pr.URL] = pr.StartTime
			}
		}
		for url, statusList := range urlStatusMap {
			mergedStatus := MergeOnlineStatus(test_result.ParseTestResults(statusList))
			entry := map[string]string{
				"time":   urlTimeMap[url],
				"online": mergedStatus.String(),
			}
			portsMap[url] = append(portsMap[url], entry)
		}
		// Clean up expired port records
		for url, history := range portsMap {
			var filteredPortHistory []map[string]string
			for _, entry := range history {
				t, err := time.Parse(time.RFC3339, entry["time"])
				if err == nil && now.Sub(t).Hours() <= float64(maxLogDays*24) {
					filteredPortHistory = append(filteredPortHistory, entry)
				}
			}
			portsMap[url] = filteredPortHistory
		}
		logData[svc.Name]["ports"] = portsMap
	}

	// Write the results to the result file
	logBytes, _ := json.MarshalIndent(logData, "", "  ")
	err := os.WriteFile(logPath, logBytes, 0644)
	if err != nil {
		log.Fatalln("Failed to write log file:", err)
	}
	return nil
}
