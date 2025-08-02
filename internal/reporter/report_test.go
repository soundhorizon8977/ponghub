package reporter

import (
	"github.com/wcy-dt/ponghub/internal/logger"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateReport tests the WriteReport function to ensure it generates a report file correctly.
func TestGenerateReport(t *testing.T) {
	logPath := "data/ponghub_log.json"
	reportPath := "data/index.html"

	logResult, err := logger.ReadLogs(logPath)
	if err != nil {
		t.Fatalf("Failed to load log data: %v", err)
	}

	err = WriteReport(logResult, reportPath)
	if err != nil {
		t.Fatalf("WriteReport failed: %v", err)
	}

	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Fatalf("Report file not generated: %s", reportPath)
	}

	f, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read generated report: %v", err)
	}
	if len(f) == 0 {
		t.Error("Generated report is empty")
	}
}

func TestMain(m *testing.M) {
	// Change the working directory to the root of the project
	root, err := filepath.Abs("../..")
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(root); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
