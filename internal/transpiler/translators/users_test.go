package translators

import (
	"testing"

	"github.com/AbHash-RixE/cloudinit-to-butane/internal/butane"
	"github.com/AbHash-RixE/cloudinit-to-butane/internal/cloudinit"
)

func TestUserTranslator_NameOnly(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{{Name: "core"}},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Passwd.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(out.Passwd.Users))
	}
	if out.Passwd.Users[0].Name != "core" {
		t.Errorf("expected name 'core', got %q", out.Passwd.Users[0].Name)
	}
}

func TestUserTranslator_WithSSHKeys(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{{
			Name:              "core",
			SSHAuthorizedKeys: []string{"ssh-ed25519 AAAA...", "ssh-rsa BBBB..."},
		}},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Passwd.Users[0].SSHAuthorizedKeys) != 2 {
		t.Errorf("expected 2 SSH keys, got %d", len(out.Passwd.Users[0].SSHAuthorizedKeys))
	}
}

func TestUserTranslator_CommaGroups(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{{
			Name:   "core",
			Groups: "sudo, docker",
		}},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	groups := out.Passwd.Users[0].Groups
	if len(groups) != 2 || groups[0] != "sudo" || groups[1] != "docker" {
		t.Errorf("expected [sudo docker], got %v", groups)
	}
}

func TestUserTranslator_WithUID(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{{
			Name: "core",
			UID:  "1000",
		}},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Passwd.Users[0].UID == nil {
		t.Fatal("expected UID to be set")
	}
	if *out.Passwd.Users[0].UID != 1000 {
		t.Errorf("expected UID 1000, got %d", *out.Passwd.Users[0].UID)
	}
}

func TestUserTranslator_InvalidUID(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{{
			Name: "core",
			UID:  "not-a-number",
		}},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Passwd.Users[0].UID != nil {
		t.Errorf("expected nil UID for invalid input, got %d", *out.Passwd.Users[0].UID)
	}
}

func TestUserTranslator_MultipleUsers(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{
			{Name: "core"},
			{Name: "admin"},
		},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Passwd.Users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(out.Passwd.Users))
	}
}

func TestUserTranslator_TopLevelGroups(t *testing.T) {
	in := &cloudinit.Config{
		Groups: []string{"wheel", "docker"},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Passwd.Groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(out.Passwd.Groups))
	}
	if out.Passwd.Groups[0].Name != "wheel" {
		t.Errorf("expected group 'wheel', got %q", out.Passwd.Groups[0].Name)
	}
}

func TestUserTranslator_Empty(t *testing.T) {
	in := &cloudinit.Config{}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(out.Passwd.Users) != 0 {
		t.Errorf("expected 0 users, got %d", len(out.Passwd.Users))
	}
	if len(out.Passwd.Groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(out.Passwd.Groups))
	}
}

func TestUserTranslator_PasswordHash(t *testing.T) {
	in := &cloudinit.Config{
		Users: []cloudinit.User{{
			Name:         "core",
			HashedPasswd: "$6$rounds=4096$abc...",
		}},
	}
	out := butane.NewDefaultConfig()

	u := NewUserTranslator()
	if err := u.Translate(in, out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Passwd.Users[0].PasswordHash != "$6$rounds=4096$abc..." {
		t.Errorf("expected password hash to be set, got %q", out.Passwd.Users[0].PasswordHash)
	}
}
