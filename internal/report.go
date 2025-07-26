package internal

import (
	"fmt"
	"github.com/wcy-dt/ponghub/protos/default_config"
	"html/template"
	"os"
)

func getLatestTime(logData LogData) string {
	var latestTime string

	for _, svcData := range logData {
		for _, entry := range svcData.ServiceHistory {
			if latestTime == "" {
				latestTime = entry.Time
			} else if entry.Time > latestTime {
				latestTime = entry.Time
			}
		}
	}

	return latestTime
}

// GenerateReport generates an HTML report from the provided log data
func GenerateReport(logData LogData, outPath string) error {
	reportEntries := logData.ParseToReportEntries()

	tmpl, err := template.New("report.html").
		Funcs(createTemplateFunc()).
		ParseFiles(default_config.GetTemplatePath())
	if err != nil {
		return fmt.Errorf("template parsing failed: %w", err)
	}

	outputFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("file creation failed: %w", err)
	}
	defer func(outputFile *os.File) {
		if err := outputFile.Close(); err != nil {
			fmt.Printf("Error closing output file: %v\n", err)
		}
	}(outputFile)

	if err := tmpl.Execute(outputFile, map[string]any{
		"Results":    reportEntries,
		"UpdateTime": getLatestTime(logData),
	}); err != nil {
		return fmt.Errorf("template execution failed: %w", err)
	}

	return nil
}

// createTemplateFunc defines custom template functions for the report
func createTemplateFunc() template.FuncMap {
	return template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b float64) float64 { return a * b },
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
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
