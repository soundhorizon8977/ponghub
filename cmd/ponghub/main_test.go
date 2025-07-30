package main

import (
	"github.com/wcy-dt/ponghub/internal/config"
	"github.com/wcy-dt/ponghub/internal/notify"
	"github.com/wcy-dt/ponghub/internal/process"
	"github.com/wcy-dt/ponghub/internal/report"
	"github.com/wcy-dt/ponghub/internal/types/default_config"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

// TestMain_append tests the main functionality when appending to an existing log file.
func TestMain_append(t *testing.T) {
	// load the default configuration
	cfg, err := config.LoadConfig(default_config.GetConfigPath())
	if err != nil {
		log.Fatalln("Error loading config at", default_config.GetConfigPath(), ":", err)
	}

	// copy log file to a temporary location for testing
	logPath := default_config.GetLogPath()
	if err := copyLogFile(logPath, tmpLogPath); err != nil {
		log.Fatalln("Error copying log file:", err)
	}

	// check services based on the configuration
	results := process.CheckServices(cfg)
	notify.NotifyResults(results)
	logData, err := process.OutputResults(results, cfg.MaxLogDays, tmpLogPath)
	if err != nil {
		log.Fatalln("Error outputting results:", err)
	}

	// generate the report based on the results
	if err := report.GenerateReport(logData, default_config.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", default_config.GetReportPath())
	}

	// Remove the temporary log file after tests
	if err := os.Remove(tmpLogPath); err != nil {
		log.Println("Error removing temporary log file:", err)
	}
}

// TestMain_new tests the main functionality when creating a new log file.
func TestMain_new(t *testing.T) {
	// load the default configuration
	cfg, err := config.LoadConfig(default_config.GetConfigPath())
	if err != nil {
		log.Fatalln("Error loading config at", default_config.GetConfigPath(), ":", err)
	}

	// check services based on the configuration
	results := process.CheckServices(cfg)
	notify.NotifyResults(results)
	logData, err := process.OutputResults(results, cfg.MaxLogDays, tmpLogPath)
	if err != nil {
		log.Fatalln("Error outputting results:", err)
	}

	// generate the report based on the results
	if err := report.GenerateReport(logData, default_config.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", default_config.GetReportPath())
	}

	// Remove the temporary log file after tests
	if err := os.Remove(tmpLogPath); err != nil {
		log.Println("Error removing temporary log file:", err)
	}
}

// copyLogFile copies the log file from srcPath to dstPath.
func copyLogFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		if err := srcFile.Close(); err != nil {
			log.Println("Error closing source file:", err)
		}
	}(srcFile)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(dstFile *os.File) {
		if err := dstFile.Close(); err != nil {
			log.Println("Error closing destination file:", err)
		}
	}(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// tmpLogPath is a temporary log file path used for testing purposes.
const tmpLogPath = "data/ponghub_log_test.json"

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
