package translators

import (
	"strings"
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

func TestCACertsTranslator_Nil(t *testing.T) {
	in := &cloudinit.Config{}
	out := butane.NewDefaultConfig()

	c := NewCACertsTranslator()
	if err := c.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(out.Storage.Files))
	}
}

func TestCACertsTranslator_SingleCert(t *testing.T) {
	in := &cloudinit.Config{
		CACerts: &cloudinit.CACerts{
			Trusted: []string{"-----BEGIN CERTIFICATE-----\nMIIC...\n-----END CERTIFICATE-----"},
		},
	}
	out := butane.NewDefaultConfig()

	c := NewCACertsTranslator()
	if err := c.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(out.Storage.Files))
	}

	file := out.Storage.Files[0]
	if file.Path != "/etc/ssl/certs/c2b-ca-0.pem" {
		t.Errorf("expected path '/etc/ssl/certs/c2b-ca-0.pem', got %q", file.Path)
	}
	if !strings.Contains(file.Contents.Inline, "BEGIN CERTIFICATE") {
		t.Error("expected certificate content")
	}
	if file.Mode == nil || *file.Mode != 420 {
		t.Error("expected mode 0644 (420)")
	}
}

func TestCACertsTranslator_MultipleCerts(t *testing.T) {
	in := &cloudinit.Config{
		CACerts: &cloudinit.CACerts{
			Trusted: []string{"cert-1", "cert-2", "cert-3"},
		},
	}
	out := butane.NewDefaultConfig()

	c := NewCACertsTranslator()
	if err := c.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(out.Storage.Files))
	}

	expectedPaths := []string{
		"/etc/ssl/certs/c2b-ca-0.pem",
		"/etc/ssl/certs/c2b-ca-1.pem",
		"/etc/ssl/certs/c2b-ca-2.pem",
	}
	for i, path := range expectedPaths {
		if out.Storage.Files[i].Path != path {
			t.Errorf("expected path %q, got %q", path, out.Storage.Files[i].Path)
		}
	}
}
