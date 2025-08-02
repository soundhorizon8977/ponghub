package configure

import (
	"github.com/wcy-dt/ponghub/internal/types/structures/configure"
	"github.com/wcy-dt/ponghub/internal/types/types/default_config"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// setDefaultConfigs sets default values for the configuration fields
func setDefaultConfigs(cfg *configure.Configure) {
	default_config.SetDefaultTimeout(&cfg.Timeout)
	default_config.SetDefaultRetry(&cfg.Retry)
	default_config.SetDefaultMaxLogDays(&cfg.MaxLogDays)

	for i := range cfg.Services {
		default_config.SetDefaultTimeout(&cfg.Services[i].Timeout)
		default_config.SetDefaultRetry(&cfg.Services[i].Retry)
	}
}

// ReadConfigs loads the configuration from a YAML file at the specified path
func ReadConfigs(path string) (*configure.Configure, error) {
	// Read the configuration file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("Error closing config file:", err)
		}
	}(f)

	// Decode the YAML configuration
	cfg := new(configure.Configure)
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		log.Fatalln("Failed to decode YAML config:", err)
	}
	// Set default values for the configuration
	setDefaultConfigs(cfg)

	if len(cfg.Services) == 0 {
		log.Fatalln("No services defined in the configuration file")
	}
	return cfg, nil
}
