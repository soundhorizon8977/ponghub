package checker

import (
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
)

// getTestResult determines the test result based on the success count and actual attempts
func getTestResult(successNum, attemptNum int) chk_result.CheckResult {
	switch successNum {
	case attemptNum:
		return chk_result.ALL
	case 0:
		return chk_result.NONE
	default:
		return chk_result.PART
	}
}
