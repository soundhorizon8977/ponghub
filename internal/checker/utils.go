package checker

import (
	"github.com/wcy-dt/ponghub/internal/types/types/test_result"
)

// getTestResult determines the test result based on the success count and actual attempts
func getTestResult(successCount, actualAttempts int) test_result.TestResult {
	switch successCount {
	case actualAttempts:
		return test_result.ALL
	case 0:
		return test_result.NONE
	default:
		return test_result.PART
	}
}
