package notify

import (
	"github.com/wcy-dt/ponghub/internal/types"
	"github.com/wcy-dt/ponghub/internal/types/default_config"
	"github.com/wcy-dt/ponghub/internal/types/test_result"
	"log"
	"os"
)

// OutputResults sends notifications based on the service check results
func OutputResults(results []types.CheckResult) {
	// find all ports with status NONE
	nonePorts := make(map[string][]string)
	for _, result := range results {
		for _, h := range result.Health {
			if h.Online == test_result.NONE {
				nonePorts[result.Name] = append(nonePorts[result.Name], h.URL)
			}
		}
		for _, a := range result.API {
			if a.Online == test_result.NONE {
				nonePorts[result.Name] = append(nonePorts[result.Name], a.URL)
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
		// if no ports are down, do nothing
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
	for service, ports := range nonePorts {
		if _, err := f.WriteString(service + "\n"); err != nil {
			log.Println("Error writing to notify file:", err)
			return
		}
		for _, port := range ports {
			if _, err := f.WriteString("\t" + port + " is unavailable.\n"); err != nil {
				log.Println("Error writing to notify file:", err)
				return
			}
		}
	}
}
