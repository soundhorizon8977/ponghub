package notifier

import (
	"github.com/wcy-dt/ponghub/internal/types/structures/checker"
	"github.com/wcy-dt/ponghub/internal/types/types/chk_result"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
	"log"
	"os"
)

// WriteNotifications sends notifications based on the service check results
func WriteNotifications(checkResult []checker.Checker) {
	// find all endpointURLs with status NONE
	nonePorts := make(map[string][]string)
	for _, serviceResult := range checkResult {
		for _, endpointResult := range serviceResult.Endpoints {
			if endpointResult.Status == chk_result.NONE {
				nonePorts[serviceResult.Name] = append(nonePorts[serviceResult.Name], endpointResult.URL)
			}
		}
	}

	// output to file default_config.GetNotifyPath()
	notifyPath := default_config.GetNotifyPath()
	if err := os.Remove(notifyPath); err != nil && !os.IsNotExist(err) {
		log.Println("Error removing notify file:", err)
		return
	}
	if len(nonePorts) == 0 {
		// if no endpointURLs are down, do nothing
		return
	}
	// new notify file
	f, err := os.Create(notifyPath)
	if err != nil {
		log.Println("Error creating notify file:", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Error closing notify file:", err)
		}
	}()
	for serviceName, endpointURLs := range nonePorts {
		if _, err := f.WriteString(serviceName + "\n"); err != nil {
			log.Println("Error writing to notify file:", err)
			return
		}
		for _, endpointURL := range endpointURLs {
			if _, err := f.WriteString("\t" + endpointURL + " is unavailable.\n"); err != nil {
				log.Println("Error writing to notify file:", err)
				return
			}
		}
	}
}
