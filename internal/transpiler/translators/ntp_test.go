package translators

import (
	"strings"
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

func TestNTPTranslator_Nil(t *testing.T) {
	in := &cloudinit.Config{}
	out := butane.NewDefaultConfig()

	n := NewNTPTranslator()
	if err := n.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Systemd.Units) != 0 {
		t.Errorf("expected 0 units, got %d", len(out.Systemd.Units))
	}
}

func TestNTPTranslator_Servers(t *testing.T) {
	in := &cloudinit.Config{
		NTP: &cloudinit.NTP{
			Servers: []string{"0.pool.ntp.org", "1.pool.ntp.org"},
		},
	}
	out := butane.NewDefaultConfig()

	n := NewNTPTranslator()
	if err := n.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Systemd.Units) != 1 {
		t.Fatalf("expected 1 unit, got %d", len(out.Systemd.Units))
	}

	unit := out.Systemd.Units[0]
	if unit.Name != "systemd-timesyncd.service" {
		t.Errorf("expected unit name 'systemd-timesyncd.service', got %q", unit.Name)
	}
	if !strings.Contains(unit.Contents, "NTP=0.pool.ntp.org") {
		t.Error("expected 'NTP=0.pool.ntp.org' in contents")
	}
	if !strings.Contains(unit.Contents, "NTP=1.pool.ntp.org") {
		t.Error("expected 'NTP=1.pool.ntp.org' in contents")
	}
}

func TestNTPTranslator_Pools(t *testing.T) {
	in := &cloudinit.Config{
		NTP: &cloudinit.NTP{
			Pools: []string{"time.example.com"},
		},
	}
	out := butane.NewDefaultConfig()

	n := NewNTPTranslator()
	if err := n.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unit := out.Systemd.Units[0]
	if !strings.Contains(unit.Contents, "Pool=time.example.com") {
		t.Error("expected 'Pool=time.example.com' in contents")
	}
}

func TestNTPTranslator_Both(t *testing.T) {
	in := &cloudinit.Config{
		NTP: &cloudinit.NTP{
			Servers: []string{"0.pool.ntp.org"},
			Pools:   []string{"time.example.com"},
		},
	}
	out := butane.NewDefaultConfig()

	n := NewNTPTranslator()
	if err := n.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unit := out.Systemd.Units[0]
	if !strings.Contains(unit.Contents, "NTP=0.pool.ntp.org") {
		t.Error("expected 'NTP=' in contents")
	}
	if !strings.Contains(unit.Contents, "Pool=time.example.com") {
		t.Error("expected 'Pool=' in contents")
	}
	if !strings.Contains(unit.Contents, "[Time]") {
		t.Error("expected '[Time]' section header")
	}
}
