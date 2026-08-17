package translators

import (
	"encoding/base64"
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

func TestFileTranslator_Basic(t *testing.T) {
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: "hello world"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(out.Storage.Files))
	}
	if out.Storage.Files[0].Path != "/etc/test" {
		t.Errorf("expected path '/etc/test', got %q", out.Storage.Files[0].Path)
	}
	if out.Storage.Files[0].Contents.Inline != "hello world" {
		t.Errorf("expected content 'hello world', got %q", out.Storage.Files[0].Contents.Inline)
	}
}

func TestFileTranslator_Permissions(t *testing.T) {
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: "x", Permissions: "0644"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Storage.Files[0].Mode == nil {
		t.Fatal("expected mode to be set")
	}
	if *out.Storage.Files[0].Mode != 420 {
		t.Errorf("expected mode 420 (0644), got %d", *out.Storage.Files[0].Mode)
	}
}

func TestFileTranslator_InvalidPermissions(t *testing.T) {
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: "x", Permissions: "not-valid"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	err := f.Translate(in, out)
	if err == nil {
		t.Fatal("expected error for invalid permissions, got nil")
	}
}

func TestFileTranslator_OwnerBoth(t *testing.T) {
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: "x", Owner: "root:root"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Storage.Files[0].User == nil || out.Storage.Files[0].User.Name != "root" {
		t.Error("expected user.name 'root'")
	}
	if out.Storage.Files[0].Group == nil || out.Storage.Files[0].Group.Name != "root" {
		t.Error("expected group.name 'root'")
	}
}

func TestFileTranslator_OwnerUserOnly(t *testing.T) {
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: "x", Owner: "core"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Storage.Files[0].User == nil || out.Storage.Files[0].User.Name != "core" {
		t.Error("expected user.name 'core'")
	}
	if out.Storage.Files[0].Group != nil {
		t.Error("expected group to be nil for owner without group")
	}
}

func TestFileTranslator_Base64(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("decoded content"))
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: encoded, Encoding: "b64"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Storage.Files[0].Contents.Inline != "decoded content" {
		t.Errorf("expected 'decoded content', got %q", out.Storage.Files[0].Contents.Inline)
	}
}

func TestFileTranslator_Base64Keyword(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("hello"))
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: encoded, Encoding: "base64"},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Storage.Files[0].Contents.Inline != "hello" {
		t.Errorf("expected 'hello', got %q", out.Storage.Files[0].Contents.Inline)
	}
}

func TestFileTranslator_Append(t *testing.T) {
	in := &cloudinit.Config{
		WriteFiles: []cloudinit.File{
			{Path: "/etc/test", Content: "appended", Append: true},
		},
	}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files[0].Append) != 1 {
		t.Fatalf("expected 1 append entry, got %d", len(out.Storage.Files[0].Append))
	}
	if out.Storage.Files[0].Append[0].Inline != "appended" {
		t.Errorf("expected 'appended', got %q", out.Storage.Files[0].Append[0].Inline)
	}
	if out.Storage.Files[0].Contents.Inline != "" {
		t.Error("expected empty contents for append mode")
	}
}

func TestFileTranslator_Empty(t *testing.T) {
	in := &cloudinit.Config{}
	out := butane.NewDefaultConfig()

	f := NewFileTranslator()
	if err := f.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Storage.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(out.Storage.Files))
	}
}
