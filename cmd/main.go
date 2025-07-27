package main

import (
	"log"

	ponghub "github.com/wcy-dt/ponghub/internal"
	"github.com/wcy-dt/ponghub/protos/default_config"
)

func main() {
	// load the default configuration
	cfg, err := ponghub.LoadConfig(default_config.GetConfigPath())
	if err != nil {
		log.Fatalln("Error loading config at", default_config.GetConfigPath(), ":", err)
	}

	// check services based on the configuration
	results := ponghub.CheckServices(cfg)
	ponghub.NotifyResults(results)
	logData, err := ponghub.OutputResults(results, cfg.MaxLogDays, default_config.GetLogPath())
	if err != nil {
		log.Fatalln("Error outputting results:", err)
	}

	// generate the report based on the results
	if err := ponghub.GenerateReport(logData, default_config.GetReportPath()); err != nil {
		log.Fatalln("Error generating report:", err)
	} else {
		log.Println("Report generated at", default_config.GetReportPath())
	}
}
