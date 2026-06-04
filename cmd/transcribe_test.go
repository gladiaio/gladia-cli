package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
	"github.com/gladiaio/gladia-cli/pkg/client/types"
)

func TestValidateOutputFormat(t *testing.T) {
	valid := []string{"text", "txt", "json", "json-full", "srt", "vtt"}
	for _, format := range valid {
		if err := validateOutputFormat(format); err != nil {
			t.Errorf("format %q: %v", format, err)
		}
	}
	if err := validateOutputFormat("table"); err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestIsHTTPURL(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"https://example.com/a.wav", true},
		{"http://example.com/a.wav", true},
		{"HTTPS://X.COM", true},
		{"/local/file.wav", false},
		{"file.wav", false},
		{"ftp://example.com", false},
	}
	for _, tc := range tests {
		if got := isHTTPURL(tc.in); got != tc.want {
			t.Errorf("isHTTPURL(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestBuildLanguageConfig(t *testing.T) {
	en := types.Language("en")
	fr := types.Language("fr")
	de := types.Language("de")

	tests := []struct {
		name          string
		langs         []types.Language
		codeSwitching bool
		codeSwitchSet bool
		wantNil       bool
		wantCodes     []string
		wantCS        bool
	}{
		{"empty", nil, false, false, true, nil, false},
		{"code switch only", nil, true, true, false, []string{}, true},
		{"code switch off explicit", nil, false, true, false, []string{}, false},
		{"single en", []types.Language{en}, false, false, false, []string{"en"}, false},
		{"single en + code switch flag", []types.Language{en}, true, true, false, []string{"en"}, true},
		{"multi languages no code switch", []types.Language{en, fr, de}, false, false, false, []string{"en", "fr", "de"}, false},
		{"multi explicit off", []types.Language{en, fr}, false, true, false, []string{"en", "fr"}, false},
		{"multi explicit on", []types.Language{en, fr}, true, true, false, []string{"en", "fr"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := buildLanguageConfig(tc.langs, tc.codeSwitching, tc.codeSwitchSet)
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantNil {
				if cfg != nil {
					t.Fatalf("expected nil config, got %+v", cfg)
				}
				return
			}
			if cfg == nil {
				t.Fatal("expected non-nil config")
			}
			if strings.Join(cfg.Languages, ",") != strings.Join(tc.wantCodes, ",") {
				t.Fatalf("languages = %v, want %v", cfg.Languages, tc.wantCodes)
			}
			if cfg.CodeSwitching != tc.wantCS {
				t.Fatalf("code_switching = %v, want %v", cfg.CodeSwitching, tc.wantCS)
			}
		})
	}
}

func TestResolveAudioSource_URL(t *testing.T) {
	client := gladia.NewGladiaClient("key", false)
	url := "https://cdn.example.com/audio.wav"
	got, err := resolveAudioSource(client, url)
	if err != nil {
		t.Fatal(err)
	}
	if got != url {
		t.Fatalf("got %q, want %q", got, url)
	}
}

func TestResolveAudioSource_missingFile(t *testing.T) {
	client := gladia.NewGladiaClient("key", false)
	_, err := resolveAudioSource(client, filepath.Join(t.TempDir(), "nope.wav"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestResolveAudioSource_upload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/upload/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"audio_url": "https://api.gladia.io/audio/123"})
	}))
	defer server.Close()

	client := gladia.NewGladiaClient("test-key", false)
	client.GladiaEndpoint = server.URL

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.wav")
	if err := os.WriteFile(path, []byte("RIFF"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := resolveAudioSource(client, path)
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://api.gladia.io/audio/123" {
		t.Fatalf("got %q", got)
	}
}

func TestTranscribeCommand_invalidOutputFormat(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "k")

	cmd := newRootCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"transcribe", "https://example.com/a.wav", "-o", "table"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid output format")
	}
}

func TestTranscribeCommand_invalidLanguage(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "k")

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"transcribe", "https://example.com/a.wav", "--language", "notalang"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid language")
	}
}

func TestTranscribeCommand_URLTextOutput(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "test-key")

	donePayload := sampleTranscriptionResult()
	donePayload.Status = "done"
	doneBody, _ := json.Marshal(donePayload)

	server := httptest.NewServer(nil)
	base := server.URL
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/transcription/":
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"result_url": base + "/result/1"})
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/result/"):
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(doneBody)
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	oldEndpoint := gladia.GladiaApiEndpoint
	gladia.GladiaApiEndpoint = server.URL
	t.Cleanup(func() { gladia.GladiaApiEndpoint = oldEndpoint })

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"transcribe", "https://example.com/audio.wav", "-o", "text"})

	out := captureStdout(t, func() {
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	})
	if got := strings.TrimSpace(out); got != "hello world" {
		t.Fatalf("stdout = %q, want hello world", got)
	}
}

func TestTranscribeCommand_codeSwitchingWithoutLanguages(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "test-key")

	var postedBody map[string]interface{}
	donePayload := sampleTranscriptionResult()
	donePayload.Status = "done"
	doneBody, _ := json.Marshal(donePayload)

	server := httptest.NewServer(nil)
	base := server.URL
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/transcription/":
			_ = json.NewDecoder(r.Body).Decode(&postedBody)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"result_url": base + "/result/1"})
		case r.Method == http.MethodGet:
			_, _ = w.Write(doneBody)
		}
	})
	defer server.Close()

	oldEndpoint := gladia.GladiaApiEndpoint
	gladia.GladiaApiEndpoint = server.URL
	t.Cleanup(func() { gladia.GladiaApiEndpoint = oldEndpoint })

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"transcribe", "https://example.com/audio.wav", "--code-switching"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	lc := postedBody["language_config"].(map[string]interface{})
	if lc["code_switching"] != true {
		t.Fatalf("code_switching = %v", lc["code_switching"])
	}
	langs, _ := lc["languages"].([]interface{})
	if len(langs) != 0 {
		t.Fatalf("languages = %v, want empty slice", lc["languages"])
	}
}

func TestTranscribeCommand_languageAndCodeSwitching(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "test-key")

	var postedBody map[string]interface{}
	donePayload := sampleTranscriptionResult()
	donePayload.Status = "done"
	doneBody, _ := json.Marshal(donePayload)

	server := httptest.NewServer(nil)
	base := server.URL
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/transcription/":
			_ = json.NewDecoder(r.Body).Decode(&postedBody)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]string{"result_url": base + "/result/1"})
		case r.Method == http.MethodGet:
			_, _ = w.Write(doneBody)
		}
	})
	defer server.Close()

	oldEndpoint := gladia.GladiaApiEndpoint
	gladia.GladiaApiEndpoint = server.URL
	t.Cleanup(func() { gladia.GladiaApiEndpoint = oldEndpoint })

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{
		"transcribe", "https://example.com/audio.wav",
		"--language", "en,fr",
		"--code-switching",
		"-o", "json",
	})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	lc, ok := postedBody["language_config"].(map[string]interface{})
	if !ok {
		t.Fatalf("language_config missing in %#v", postedBody)
	}
	if lc["code_switching"] != true {
		t.Fatalf("code_switching = %v, want true when --code-switching is set", lc["code_switching"])
	}
	langs, ok := lc["languages"].([]interface{})
	if !ok || len(langs) != 2 {
		t.Fatalf("languages = %v", lc["languages"])
	}
}
