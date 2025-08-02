package main

import (
	"github.com/wcy-dt/ponghub/internal/checker"
	"github.com/wcy-dt/ponghub/internal/configure"
	"github.com/wcy-dt/ponghub/internal/logger"
	"github.com/wcy-dt/ponghub/internal/notifier"
	"github.com/wcy-dt/ponghub/internal/reporter"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
	"log"
)

func main() {
	// load the default configuration
	cfg, err := configure.ReadConfigs(default_config.GetConfigPath())
	if err != nil {
		log.Fatalln("Error loading config at", default_config.GetConfigPath(), ":", err)
	}

	// check services based on the configuration
	results := checker.CheckServices(cfg)
	notifier.WriteNotifications(results)
	logData, err := logger.OutputResults(results, cfg.MaxLogDays, default_config.GetLogPath())
	if err != nil {
		log.Fatalln("Error outputting results:", err)
	}

	// generate the report based on the results
	if err := reporter.GenerateReport(logData, default_config.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", default_config.GetReportPath())
	}
}
