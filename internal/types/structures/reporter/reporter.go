package reporter

// Data structures for logging and reporting
type (
	// HistoryEntry represents a single history entry
	HistoryEntry struct {
		Time         string `json:"time"`
		Status       string `json:"online"`
		ResponseTime int    `json:"response_time,omitempty"`
	}

	History []HistoryEntry
	Port    map[string]History

	// ReportEntry represents the result of checking a service
	ReportEntry struct {
		Name         string
		History      History
		Ports        Port
		Availability float64
	}
)
