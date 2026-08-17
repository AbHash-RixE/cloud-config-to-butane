package cloudinit

import (
	"strings"
	"testing"
)

func TestParse_WithHeader(t *testing.T) {
	input := []byte(`#cloud-config

users:
  - name: core
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Users) != 1 || cfg.Users[0].Name != "core" {
		t.Errorf("expected user 'core', got %+v", cfg.Users)
	}
}

func TestParse_WithJinjaHeader(t *testing.T) {
	input := []byte(`## template: jinja
#cloud-config

users:
  - name: core
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Users) != 1 || cfg.Users[0].Name != "core" {
		t.Errorf("expected user 'core', got %+v", cfg.Users)
	}
}

func TestParse_NoHeader(t *testing.T) {
	input := []byte(`users:
  - name: core
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Users) != 1 || cfg.Users[0].Name != "core" {
		t.Errorf("expected user 'core', got %+v", cfg.Users)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	input := []byte(`{{{invalid`)
	_, err := Parse(input)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestParse_EmptyInput(t *testing.T) {
	input := []byte(``)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Users) != 0 {
		t.Errorf("expected no users, got %d", len(cfg.Users))
	}
}

func TestParse_RuncmdString(t *testing.T) {
	input := []byte(`runcmd:
  - "echo hi"
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.RunCmd) != 1 {
		t.Fatalf("expected 1 runcmd, got %d", len(cfg.RunCmd))
	}
	if string(cfg.RunCmd[0]) != "echo hi" {
		t.Errorf("expected 'echo hi', got %q", string(cfg.RunCmd[0]))
	}
}

func TestParse_RuncmdList(t *testing.T) {
	input := []byte(`runcmd:
  - [echo, hello, world]
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.RunCmd) != 1 {
		t.Fatalf("expected 1 runcmd, got %d", len(cfg.RunCmd))
	}
	cmd := string(cfg.RunCmd[0])
	if !strings.Contains(cmd, "echo") || !strings.Contains(cmd, "hello") || !strings.Contains(cmd, "world") {
		t.Errorf("expected shell-quoted args, got %q", cmd)
	}
}

func TestParse_BootCmd(t *testing.T) {
	input := []byte(`bootcmd:
  - "echo test"
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.BootCmd) != 1 {
		t.Fatalf("expected 1 bootcmd, got %d", len(cfg.BootCmd))
	}
	if string(cfg.BootCmd[0]) != "echo test" {
		t.Errorf("expected 'echo test', got %q", string(cfg.BootCmd[0]))
	}
}

func TestParse_NTP(t *testing.T) {
	input := []byte(`ntp:
  servers:
    - 0.pool.ntp.org
  pools:
    - time.example.com
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.NTP == nil {
		t.Fatal("expected NTP config, got nil")
	}
	if len(cfg.NTP.Servers) != 1 || cfg.NTP.Servers[0] != "0.pool.ntp.org" {
		t.Errorf("expected server '0.pool.ntp.org', got %v", cfg.NTP.Servers)
	}
	if len(cfg.NTP.Pools) != 1 || cfg.NTP.Pools[0] != "time.example.com" {
		t.Errorf("expected pool 'time.example.com', got %v", cfg.NTP.Pools)
	}
}

func TestParse_CACerts(t *testing.T) {
	input := []byte(`ca_certs:
  trusted:
    - |
      -----BEGIN CERTIFICATE-----
      MIIC...
      -----END CERTIFICATE-----
`)
	cfg, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CACerts == nil {
		t.Fatal("expected CACerts config, got nil")
	}
	if len(cfg.CACerts.Trusted) != 1 {
		t.Fatalf("expected 1 trusted cert, got %d", len(cfg.CACerts.Trusted))
	}
	if !strings.Contains(cfg.CACerts.Trusted[0], "BEGIN CERTIFICATE") {
		t.Errorf("expected certificate content, got %q", cfg.CACerts.Trusted[0])
	}
}
