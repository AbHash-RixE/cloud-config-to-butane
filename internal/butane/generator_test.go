package butane

import (
	"strings"
	"testing"
)

func TestGenerate_Defaults(t *testing.T) {
	cfg := NewDefaultConfig()
	data, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "variant: flatcar") {
		t.Error("expected 'variant: flatcar' in output")
	}
	if !strings.Contains(out, "version: 1.1.0") {
		t.Error("expected 'version: 1.1.0' in output")
	}
}

func TestGenerate_EmptyConfig(t *testing.T) {
	cfg := NewDefaultConfig()
	data, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestGenerate_WithUser(t *testing.T) {
	cfg := NewDefaultConfig()
	cfg.Passwd.Users = []User{
		{Name: "core"},
	}
	data, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "name: core") {
		t.Error("expected 'name: core' in output")
	}
}

func TestGenerate_WithFile(t *testing.T) {
	cfg := NewDefaultConfig()
	mode := 420
	cfg.Storage.Files = []File{
		{
			Path:     "/etc/test",
			Mode:     &mode,
			Contents: Content{Inline: "hello"},
		},
	}
	data, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "path: /etc/test") {
		t.Error("expected file path in output")
	}
	if !strings.Contains(out, "mode: 420") {
		t.Error("expected mode in output")
	}
}

func TestGenerate_WithUnit(t *testing.T) {
	cfg := NewDefaultConfig()
	enabled := true
	cfg.Systemd.Units = []Unit{
		{
			Name:     "test.service",
			Enabled:  &enabled,
			Contents: "[Service]\nExecStart=/bin/true",
		},
	}
	data, err := Generate(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := string(data)
	if !strings.Contains(out, "name: test.service") {
		t.Error("expected unit name in output")
	}
}
