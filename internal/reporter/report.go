package reporter

import (
	"fmt"
	"github.com/wcy-dt/ponghub/internal/types/structures/logger"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
	"html/template"
	"os"
)

// WriteReport generates an HTML report from the provided log data
func WriteReport(logResult logger.Logger, reportPath string) error {
	reportResult := logResult.ParseToReportEntries()

	tmpl, err := template.New("report.html").
		Funcs(createTemplateFunc()).
		ParseFiles(default_config.GetTemplatePath())
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	reportFile, err := os.Create(reportPath)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer func(reportFile *os.File) {
		if err := reportFile.Close(); err != nil {
			fmt.Printf("Error closing output file: %v\n", err)
		}
	}(reportFile)

	if err := tmpl.Execute(reportFile, map[string]any{
		"ReportResult": reportResult,
		"UpdateTime":   getLatestTime(logResult),
	}); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// getLatestTime retrieves the latest time from the log data
func getLatestTime(logResult logger.Logger) string {
	var latestTime string

	for _, serviceResult := range logResult {
		for _, serviceHistoryEntry := range serviceResult.ServiceHistory {
			if latestTime == "" {
				latestTime = serviceHistoryEntry.Time
			} else if serviceHistoryEntry.Time > latestTime {
				latestTime = serviceHistoryEntry.Time
			}
		}
	}

	return latestTime
}

// createTemplateFunc defines custom template functions for the report
func createTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"div": func(a, b int) float64 {
			if b == 0 {
				return 0
			}
			return float64(a) / float64(b)
		},
		"until": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
		},
	}
}
