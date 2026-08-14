package cloudinit

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Config, error) {
	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse cloud-config YAML: %w", err)
	}

	return &cfg, nil
}

// shellQuote safety: escaping inner single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

func (r *RuncmdItem) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		*r = RuncmdItem(value.Value)
		return nil
	}

	if value.Kind == yaml.SequenceNode {
		var parts []string
		for _, node := range value.Content {

			// Validate for strings only
			if node.Kind != yaml.ScalarNode {
				return fmt.Errorf("runcmd list arguments must be scalars, got node kind: %v", node.Kind)
			}

			parts = append(parts, shellQuote(node.Value))
		}
		*r = RuncmdItem(strings.Join(parts, " "))
		return nil
	}

	return fmt.Errorf("runcmd item must be a string or a list of strings, got node kind: %v", value.Kind)
}
