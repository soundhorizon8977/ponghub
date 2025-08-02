package main

import (
	"github.com/wcy-dt/ponghub/internal/config"
	"github.com/wcy-dt/ponghub/internal/notify"
	"github.com/wcy-dt/ponghub/internal/process"
	"github.com/wcy-dt/ponghub/internal/report"
	"github.com/wcy-dt/ponghub/internal/types/default_config"
	"log"
)

func main() {
	// load the default configuration
	cfg, err := config.LoadConfig(default_config.GetConfigPath())
	if err != nil {
		log.Fatalln("Error loading config at", default_config.GetConfigPath(), ":", err)
	}

	// check services based on the configuration
	results := process.CheckServices(cfg)
	notify.OutputResults(results)
	logData, err := process.OutputResults(results, cfg.MaxLogDays, default_config.GetLogPath())
	if err != nil {
		log.Fatalln("Error outputting results:", err)
	}

	// generate the report based on the results
	if err := report.GenerateReport(logData, default_config.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", default_config.GetReportPath())
	}
}
