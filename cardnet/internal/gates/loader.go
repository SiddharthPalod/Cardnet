package gates

import (
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads gate configuration from YAML file
func LoadConfig(path string) (Config, error) {
	if path == "" {
		path = "configs/gates.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return cfg, nil
}
