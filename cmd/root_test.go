package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCommand_help(t *testing.T) {
	cmd := newRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, sub := range []string{"transcribe", "auth", "languages"} {
		if !strings.Contains(out, sub) {
			t.Fatalf("help missing %q:\n%s", sub, out)
		}
	}
}

func TestAuthSet_writesConfig(t *testing.T) {
	home := withTempHome(t)
	t.Setenv(envGladiaAPIKey, "")

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"auth", "set", "my-secret"})

	captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})

	key, err := GetGladiaKeyFromFile()
	if err != nil {
		t.Fatal(err)
	}
	if key != "my-secret" {
		t.Fatalf("got %q", key)
	}

	info, err := os.Stat(filepath.Join(home, configFilename))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o", info.Mode().Perm())
	}
}

func TestLanguagesCommand(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"languages"})
	out := captureStdout(t, func() {
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "en:") {
		t.Fatalf("expected language listing, got %q", out)
	}
	if strings.Contains(out, "Error parsing") {
		t.Fatalf("language display errors in output: %q", out)
	}
	if !strings.Contains(out, "at: Asturian") {
		t.Fatalf("expected Asturian entry, got %q", out)
	}
	if !strings.Contains(out, "jp: Japanese") {
		t.Fatalf("expected Japanese entry, got %q", out)
	}
	if !strings.Contains(out, "mymr: Burmese") {
		t.Fatalf("expected Burmese entry, got %q", out)
	}
}
