package translators

import (
	"strings"
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

func TestBootCmdTranslator_Empty(t *testing.T) {
	in := &cloudinit.Config{}
	out := butane.NewDefaultConfig()

	b := NewBootCmdTranslator()
	if err := b.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(out.Storage.Files))
	}
	if len(out.Systemd.Units) != 0 {
		t.Errorf("expected 0 units, got %d", len(out.Systemd.Units))
	}
}

func TestBootCmdTranslator_Basic(t *testing.T) {
	in := &cloudinit.Config{
		BootCmd: []cloudinit.RuncmdItem{"echo boot"},
	}
	out := butane.NewDefaultConfig()

	b := NewBootCmdTranslator()
	if err := b.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 1 {
		t.Fatalf("expected 1 script file, got %d", len(out.Storage.Files))
	}
	if out.Storage.Files[0].Path != "/opt/c2b/bootcmd.sh" {
		t.Errorf("expected path '/opt/c2b/bootcmd.sh', got %q", out.Storage.Files[0].Path)
	}

	content := out.Storage.Files[0].Contents.Inline
	if !strings.Contains(content, "echo boot") {
		t.Error("expected 'echo boot' in script")
	}
}

func TestBootCmdTranslator_UnitContents(t *testing.T) {
	in := &cloudinit.Config{
		BootCmd: []cloudinit.RuncmdItem{"echo boot"},
	}
	out := butane.NewDefaultConfig()

	b := NewBootCmdTranslator()
	if err := b.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Systemd.Units) != 1 {
		t.Fatalf("expected 1 unit, got %d", len(out.Systemd.Units))
	}

	unit := out.Systemd.Units[0]
	if unit.Name != "c2b-bootcmd.service" {
		t.Errorf("expected unit name 'c2b-bootcmd.service', got %q", unit.Name)
	}
	if unit.Enabled == nil || !*unit.Enabled {
		t.Error("expected unit to be enabled")
	}
	if !strings.Contains(unit.Contents, "Before=network.target") {
		t.Error("expected 'Before=network.target' in unit contents")
	}
	if !strings.Contains(unit.Contents, "DefaultDependencies=no") {
		t.Error("expected 'DefaultDependencies=no' in unit contents")
	}
}

func TestBootCmdTranslator_MultipleCommands(t *testing.T) {
	in := &cloudinit.Config{
		BootCmd: []cloudinit.RuncmdItem{
			"echo first",
			"echo second",
		},
	}
	out := butane.NewDefaultConfig()

	b := NewBootCmdTranslator()
	if err := b.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := out.Storage.Files[0].Contents.Inline
	if !strings.Contains(content, "echo first") || !strings.Contains(content, "echo second") {
		t.Error("expected both commands in script")
	}
}
