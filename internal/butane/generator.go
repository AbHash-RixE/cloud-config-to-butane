package butane

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Generate(cfg *Config) ([]byte, error) {
	// Safety net
	if cfg.Variant == "" {
		cfg.Variant = "flatcar"
	}
	if cfg.Version == "" {
		cfg.Version = "1.1.0"
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to generate butane YAML: %w", err)
	}

	return data, nil
}
