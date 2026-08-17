package translators

import (
	"strings"
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

func TestRuncmdTranslator_SingleString(t *testing.T) {
	in := &cloudinit.Config{
		RunCmd: []cloudinit.RuncmdItem{"echo hello"},
	}
	out := butane.NewDefaultConfig()

	r := NewRunCmdTranslator()
	if err := r.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 1 {
		t.Fatalf("expected 1 script file, got %d", len(out.Storage.Files))
	}
	if out.Storage.Files[0].Path != "/opt/c2b/runcmd.sh" {
		t.Errorf("expected path '/opt/c2b/runcmd.sh', got %q", out.Storage.Files[0].Path)
	}
	content := out.Storage.Files[0].Contents.Inline
	if !strings.Contains(content, "set -e") {
		t.Error("expected 'set -e' in script")
	}
	if !strings.Contains(content, "echo hello") {
		t.Error("expected 'echo hello' in script")
	}
}

func TestRuncmdTranslator_MultipleCmds(t *testing.T) {
	in := &cloudinit.Config{
		RunCmd: []cloudinit.RuncmdItem{
			"echo one",
			"echo two",
			"echo three",
		},
	}
	out := butane.NewDefaultConfig()

	r := NewRunCmdTranslator()
	if err := r.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := out.Storage.Files[0].Contents.Inline
	if !strings.Contains(content, "echo one") || !strings.Contains(content, "echo two") || !strings.Contains(content, "echo three") {
		t.Error("expected all 3 commands in script")
	}
}

func TestRuncmdTranslator_Empty(t *testing.T) {
	in := &cloudinit.Config{}
	out := butane.NewDefaultConfig()

	r := NewRunCmdTranslator()
	if err := r.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 0 {
		t.Errorf("expected 0 files for empty runcmd, got %d", len(out.Storage.Files))
	}
	if len(out.Systemd.Units) != 0 {
		t.Errorf("expected 0 units for empty runcmd, got %d", len(out.Systemd.Units))
	}
}

func TestRuncmdTranslator_SystemdUnit(t *testing.T) {
	in := &cloudinit.Config{
		RunCmd: []cloudinit.RuncmdItem{"echo test"},
	}
	out := butane.NewDefaultConfig()

	r := NewRunCmdTranslator()
	if err := r.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Systemd.Units) != 1 {
		t.Fatalf("expected 1 systemd unit, got %d", len(out.Systemd.Units))
	}
	unit := out.Systemd.Units[0]
	if unit.Name != "c2b-runcmd.service" {
		t.Errorf("expected unit name 'c2b-runcmd.service', got %q", unit.Name)
	}
	if unit.Enabled == nil || !*unit.Enabled {
		t.Error("expected unit to be enabled")
	}
	if !strings.Contains(unit.Contents, "oneshot") {
		t.Error("expected oneshot service type")
	}
}
