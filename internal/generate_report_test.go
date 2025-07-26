package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// Change the working directory to the root of the project
	root, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(root); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestGenerateReport(t *testing.T) {
	logPath := "data/ponghub_log.json"
	outPath := "data/index.html"

	err := GenerateReport(logPath, outPath)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatalf("Report file not generated: %s", outPath)
	}

	f, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Failed to read generated report: %v", err)
	}
	if len(f) == 0 {
		t.Error("Generated report is empty")
	}
}
