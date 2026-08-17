package transpiler

import (
	"strings"
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/transpiler/translators"
)

func TestEngine_FullPipeline(t *testing.T) {
	input := []byte(`#cloud-config

users:
  - name: core
    ssh_authorized_keys:
      - ssh-ed25519 AAAA...
    groups: sudo
    uid: "1000"

write_files:
  - path: /etc/kubernetes/kubelet.conf
    owner: root:root
    permissions: "0644"
    content: |
      apiVersion: v1
      kind: Config

bootcmd:
  - "echo boot"

ntp:
  servers:
    - 0.pool.ntp.org

runcmd:
  - "echo hello"
`)

	cloudCfg, err := cloudinit.Parse(input)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	engine := NewEngine(
		translators.NewBootCmdTranslator(),
		translators.NewUserTranslator(),
		translators.NewFileTranslator(),
		translators.NewRunCmdTranslator(),
		translators.NewNTPTranslator(),
		translators.NewCACertsTranslator(),
	)

	butaneCfg, err := engine.Run(cloudCfg)
	if err != nil {
		t.Fatalf("transpile failed: %v", err)
	}

	if butaneCfg.Variant != "flatcar" {
		t.Errorf("expected variant 'flatcar', got %q", butaneCfg.Variant)
	}
	if butaneCfg.Version != "1.1.0" {
		t.Errorf("expected version '1.1.0', got %q", butaneCfg.Version)
	}
	if len(butaneCfg.Passwd.Users) != 1 {
		t.Errorf("expected 1 user, got %d", len(butaneCfg.Passwd.Users))
	}
	if butaneCfg.Passwd.Users[0].Name != "core" {
		t.Errorf("expected user 'core', got %q", butaneCfg.Passwd.Users[0].Name)
	}
	if butaneCfg.Passwd.Users[0].UID == nil || *butaneCfg.Passwd.Users[0].UID != 1000 {
		t.Errorf("expected UID 1000")
	}

	fileCount := 0
	for _, f := range butaneCfg.Storage.Files {
		if f.Path == "/etc/kubernetes/kubelet.conf" {
			fileCount++
			if f.Group == nil || f.Group.Name != "root" {
				t.Error("expected file group 'root'")
			}
		}
		if f.Path == "/opt/c2b/bootcmd.sh" {
			fileCount++
		}
		if f.Path == "/opt/c2b/runcmd.sh" {
			fileCount++
		}
	}
	if fileCount != 3 {
		t.Errorf("expected 3 specific files, got %d", fileCount)
	}

	hasNTP := false
	for _, u := range butaneCfg.Systemd.Units {
		if u.Name == "systemd-timesyncd.service" {
			hasNTP = true
		}
	}
	if !hasNTP {
		t.Error("expected systemd-timesyncd unit")
	}
}

func TestEngine_EmptyConfig(t *testing.T) {
	cloudCfg := &cloudinit.Config{}

	engine := NewEngine(
		translators.NewBootCmdTranslator(),
		translators.NewUserTranslator(),
		translators.NewFileTranslator(),
		translators.NewRunCmdTranslator(),
		translators.NewNTPTranslator(),
		translators.NewCACertsTranslator(),
	)

	butaneCfg, err := engine.Run(cloudCfg)
	if err != nil {
		t.Fatalf("transpile failed: %v", err)
	}

	if butaneCfg.Variant != "flatcar" {
		t.Errorf("expected variant 'flatcar', got %q", butaneCfg.Variant)
	}
	if butaneCfg.Version != "1.1.0" {
		t.Errorf("expected version '1.1.0', got %q", butaneCfg.Version)
	}
}

func TestEngine_TranslatorErrorPropagates(t *testing.T) {
	cloudCfg := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/bad", Permissions: "not-octal"},
		},
	}

	engine := NewEngine(
		translators.NewFileTranslator(),
	)

	_, err := engine.Run(cloudCfg)
	if err == nil {
		t.Fatal("expected error from invalid permissions, got nil")
	}
	if !strings.Contains(err.Error(), "write_files") {
		t.Errorf("expected error to mention module name, got: %v", err)
	}
}
