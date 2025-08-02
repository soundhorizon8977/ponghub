package logger

import (
	"github.com/wcy-dt/ponghub/internal/types/structures/reporter"
	"github.com/wcy-dt/ponghub/internal/types/types/test_result"
	"log"
	"time"
)

// CleanExpiredEntries removes entries older than maxDays from the history entry list.
func (h *History) CleanExpiredEntries(maxDays int) {
	if maxDays <= 0 {
		log.Println("Max days for cleaning history is not set or invalid, skipping cleaning.")
		return // No cleaning needed if maxDays is not set
	}

	cutoffTime := time.Now().AddDate(0, 0, -maxDays)
	var cleanedHistory History

	for _, entry := range *h {
		entryTime, err := time.Parse(time.RFC3339, entry.Time)
		if err != nil {
			log.Printf("Error parsing time %s: %v", entry.Time, err)
			continue // Skip entries with invalid time format
		}
		if entryTime.After(cutoffTime) {
			cleanedHistory = append(cleanedHistory, entry)
		}
	}

	*h = cleanedHistory
}

// AddEntry adds a new entry to the history entry list.
func (h *History) AddEntry(entry HistoryEntry) {
	*h = append(*h, entry)
}

// ParseToReportEntries converts LogData to a slice of ReportEntry, calculating availability for each service.
func (l Logger) ParseToReportEntries() []reporter.ReportEntry {
	var reportEntries []reporter.ReportEntry
	for svcName, svcData := range l {
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

		reportEntries = append(reportEntries, reporter.ReportEntry{
			Name:         svcName,
			History:      svcData.ServiceHistory.loggerHistoryToReporterHistory(),
			Ports:        svcData.PortsData.loggerPortToReporterPort(),
			Availability: availability,
		})
	}
	return reportEntries
}

func (h History) loggerHistoryToReporterHistory() reporter.History {
	var reporterHistory reporter.History
	for _, entry := range h {
		reporterHistory = append(reporterHistory, reporter.HistoryEntry{
			Time:         entry.Time,
			Status:       entry.Status,
			ResponseTime: entry.ResponseTime,
		})
	}
	return reporterHistory
}

func (p Port) loggerPortToReporterPort() reporter.Port {
	reporterPort := make(reporter.Port)
	for port, history := range p {
		reporterHistory := history.loggerHistoryToReporterHistory()
		reporterPort[port] = reporterHistory
	}
	return reporterPort
}
