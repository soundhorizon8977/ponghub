package reporter

import (
	"github.com/wcy-dt/ponghub/internal/types/structures/logger"
	"log"
)

// ParseLogResult converts logger.Logger data into a reporter.Reporter format.
func ParseLogResult(logResult logger.Logger) Reporter {
	report := make(Reporter)
	for serviceName, serviceLog := range logResult {
		if len(serviceLog.ServiceHistory) == 0 {
			log.Printf("No history data for service %s", serviceName)
			continue // Skip services with no history data
		}

		// Convert logger.Endpoints to reporter.Endpoints
		endpoints := make(Endpoints)
		for url, endpointLog := range serviceLog.Endpoints {
			var endpointHistory History
			for _, entry := range endpointLog {
				endpointHistory = append(endpointHistory, HistoryEntry{
					Time:         entry.Time,
					Status:       entry.Status,
					ResponseTime: entry.ResponseTime,
				})
			}
			endpoints[url] = Endpoint{
				EndpointHistory: endpointHistory,
			}
		}

		// convert logger.ServiceHistory to reporter.ServiceHistory
		serviceHistory := make(History, len(serviceLog.ServiceHistory))
		for i, entry := range serviceLog.ServiceHistory {
			serviceHistory[i] = HistoryEntry{
				Time:         entry.Time,
				Status:       entry.Status,
				ResponseTime: entry.ResponseTime,
			}
		}

		newService := Service{
			ServiceHistory: serviceHistory,
			Endpoints:      endpoints,
		}
		report[serviceName] = newService
	}
	return report
}
