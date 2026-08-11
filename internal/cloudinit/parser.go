package cloudinit

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse cloud-config YAML: %w", err)
	}

	return &cfg, nil
}
