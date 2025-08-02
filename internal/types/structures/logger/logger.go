package logger

type (
	// HistoryEntry represents a single history entry
	HistoryEntry struct {
		Time         string `json:"time"`
		Status       string `json:"online"`
		ResponseTime int    `json:"response_time,omitempty"`
	}

	History []HistoryEntry
	Port    map[string]History

	// Entry represents log data for a service
	Entry struct {
		ServiceHistory History `json:"service_history"`
		PortsData      Port    `json:"ports"`
	}

	// Logger represents the entire log structure
	Logger map[string]Entry
)
