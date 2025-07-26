package internal

import (
	"github.com/wcy-dt/ponghub/protos/default_config"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// SetDefaultFields sets default values for the configuration fields
func SetDefaultFields(cfg *Config) {
	default_config.SetDefaultTimeout(&cfg.Timeout)
	default_config.SetDefaultRetry(&cfg.Retry)
	default_config.SetDefaultMaxLogDays(&cfg.MaxLogDays)

	for i := range cfg.Services {
		default_config.SetDefaultTimeout(&cfg.Services[i].Timeout)
		default_config.SetDefaultRetry(&cfg.Services[i].Retry)
	}
}

// LoadConfig loads the configuration from a YAML file at the specified path
func LoadConfig(path string) (*Config, error) {
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
	cfg := new(Config)
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		log.Fatalln("Failed to decode YAML config:", err)
	}
	// Set default values for the configuration
	SetDefaultFields(cfg)

	if len(cfg.Services) == 0 {
		log.Fatalln("No services defined in the configuration file")
	}
	return cfg, nil
}
