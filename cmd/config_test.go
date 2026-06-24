package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveAPIKey_fromEnv(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "env-key")
	t.Cleanup(func() { t.Setenv(envGladiaAPIKey, "") })

	key, err := ResolveAPIKey("")
	if err != nil {
		t.Fatal(err)
	}
	if key != "env-key" {
		t.Fatalf("got %q, want env-key", key)
	}
}

func TestResolveAPIKey_fromFile(t *testing.T) {
	home := withTempHome(t)
	t.Setenv(envGladiaAPIKey, "")
	path := filepath.Join(home, configFilename)
	if err := os.WriteFile(path, []byte("file-key\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	key, err := ResolveAPIKey("")
	if err != nil {
		t.Fatal(err)
	}
	if key != "file-key" {
		t.Fatalf("got %q, want file-key", key)
	}
}

func TestResolveAPIKey_fromFlag(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "")

	key, err := ResolveAPIKey("flag-key")
	if err != nil {
		t.Fatal(err)
	}
	if key != "flag-key" {
		t.Fatalf("got %q, want flag-key", key)
	}
}

func TestResolveAPIKey_priorityEnvOverFileAndFlag(t *testing.T) {
	home := withTempHome(t)
	t.Setenv(envGladiaAPIKey, "env-wins")
	path := filepath.Join(home, configFilename)
	if err := os.WriteFile(path, []byte("file-key\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	key, err := ResolveAPIKey("flag-key")
	if err != nil {
		t.Fatal(err)
	}
	if key != "env-wins" {
		t.Fatalf("got %q, want env-wins", key)
	}
}

func TestResolveAPIKey_priorityFileOverFlag(t *testing.T) {
	home := withTempHome(t)
	t.Setenv(envGladiaAPIKey, "")
	path := filepath.Join(home, configFilename)
	if err := os.WriteFile(path, []byte("file-wins\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	key, err := ResolveAPIKey("flag-key")
	if err != nil {
		t.Fatal(err)
	}
	if key != "file-wins" {
		t.Fatalf("got %q, want file-wins", key)
	}
}

func TestResolveAPIKey_missing(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "")

	_, err := ResolveAPIKey("")
	if err == nil {
		t.Fatal("expected error when no API key configured")
	}
	msg := err.Error()
	for _, part := range []string{"GLADIA_API_KEY", "gladia auth set", "app.gladia.io"} {
		if !strings.Contains(msg, part) {
			t.Fatalf("error %q missing %q", msg, part)
		}
	}
}

func TestSaveGladiaKeyToFile_contentAndPermissions(t *testing.T) {
	home := withTempHome(t)

	captureStdout(t, func() {
		if err := SaveGladiaKeyToFile("  secret-key  \n"); err != nil {
			t.Fatal(err)
		}
	})

	path := filepath.Join(home, configFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "secret-key\n" {
		t.Fatalf("file content = %q, want %q", data, "secret-key\n")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %o, want 0600", info.Mode().Perm())
	}
}

func TestGetGladiaKeyFromFile_missing(t *testing.T) {
	withTempHome(t)
	_, err := GetGladiaKeyFromFile()
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}
