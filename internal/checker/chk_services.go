package checker

import (
	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
	"github.com/wcy-dt/ponghub/internal/types/types/test_result"
	"time"
)

// CheckServices checks all services defined in the configuration
func CheckServices(cfg *configure.Configure) []checker.Checker {
	var results []checker.Checker
	for _, svc := range cfg.Services {
		// start timer
		svcStart := time.Now()

		totalAttempts := 0
		successCount := 0
		totalPorts := 0
		onlinePorts := 0

		// check health ports
		var healthResults []checker.Port
		for _, h := range svc.Health {
			pr := checkPort(&h, svc.Timeout, svc.Retry, svc.Name)
			healthResults = append(healthResults, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == test_result.ALL {
				onlinePorts++
			}
		}

		// check API ports
		var apiResults []checker.Port
		for _, a := range svc.API {
			pr := checkPort(&a, svc.Timeout, svc.Retry, svc.Name)
			apiResults = append(apiResults, pr)
			totalAttempts += pr.TotalAttempts
			successCount += pr.SuccessCount
			totalPorts++
			if pr.Online == test_result.ALL {
				onlinePorts++
			}
		}

		// end timer
		svcEnd := time.Now()

		res := checker.Checker{
			Name:          svc.Name,
			Online:        getTestResult(onlinePorts, totalPorts),
			Health:        healthResults,
			API:           apiResults,
			StartTime:     svcStart.Format(time.RFC3339),
			EndTime:       svcEnd.Format(time.RFC3339),
			TotalAttempts: totalAttempts,
			SuccessCount:  successCount,
		}
		results = append(results, res)
	}
	return results
}
