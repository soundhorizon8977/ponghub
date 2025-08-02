package reporter

// Data structures for logging and reporting
type (
	// HistoryEntry represents a single history entry
	HistoryEntry struct {
		Time         string
		Status       string
		ResponseTime int
	}

	History   []HistoryEntry
	Endpoints map[string]History

	// Reporter represents the result of checking a service
	Reporter struct {
		Name           string
		ServiceHistory History
		Endpoints      Endpoints
		Availability   float64
	}
)
