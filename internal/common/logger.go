package common

import (
	"encoding/json"
	"github.com/wcy-dt/ponghub/internal/types/structures/logger"
	"os"
)

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
