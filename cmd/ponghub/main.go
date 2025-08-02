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
	checkResult := checker.CheckServices(cfg)
	notifier.WriteNotifications(checkResult)
	logResult, err := logger.GetLogs(checkResult, cfg.MaxLogDays, default_config.GetLogPath())
	if err != nil {
		log.Fatalln("Error outputting checkResult:", err)
	}

	// generate the report based on the checkResult
	if err := reporter.WriteReport(logResult, default_config.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", default_config.GetReportPath())
	}
}
