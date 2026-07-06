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

func TestValidateModel(t *testing.T) {
	for _, model := range []string{"", "solaria-1", "solaria-3", "solaria 3", " Solaria-3 "} {
		if err := validateModel(model); err != nil {
			t.Errorf("model %q: %v", model, err)
		}
	}
	if err := validateModel("solaria-2"); err == nil {
		t.Fatal("expected error for unknown model")
	}
}

func TestValidateModelConfig_solaria3(t *testing.T) {
	en := types.LanguageEn
	fr := types.LanguageFr
	ja := types.LanguageJp

	if err := validateModelConfig("solaria-3", []types.Language{en}, false, false); err != nil {
		t.Fatalf("single supported language: %v", err)
	}
	if err := validateModelConfig("solaria-3", nil, false, false); err != nil {
		t.Fatalf("no language should be allowed: %v", err)
	}
	if err := validateModelConfig("solaria-3", []types.Language{en, fr}, false, false); err == nil {
		t.Fatal("expected error for multiple languages")
	} else if !strings.Contains(err.Error(), "only one language") || !strings.Contains(err.Error(), "en, fr") {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateModelConfig("solaria-3", []types.Language{ja}, false, false); err == nil {
		t.Fatal("expected error for unsupported language")
	}
	if err := validateModelConfig("solaria-3", []types.Language{en}, true, true); err == nil {
		t.Fatal("expected error for code switching")
	}
	if err := validateModelConfig("solaria-1", nil, false, false); err != nil {
		t.Fatalf("solaria-1 should not require language: %v", err)
	}
}

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

func TestTranscribeCommand_invalidModel(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "k")

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"transcribe", "https://example.com/a.wav", "--model", "solaria-2"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid model")
	}
}

func TestTranscribeCommand_invalidSolaria3MultipleLanguages(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "k")

	cmd := newRootCmd()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"transcribe", "https://example.com/a.wav", "--model", "solaria-3", "--language", "en,fr"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for solaria-3 with multiple languages")
	}
	if !strings.Contains(err.Error(), "only one language") || !strings.Contains(err.Error(), "en, fr") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTranscribeCommand_solaria3WithoutLanguage(t *testing.T) {
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
		case r.Method == http.MethodPost && r.URL.Path == "/v2/pre-recorded":
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
	cmd.SetArgs([]string{"transcribe", "https://example.com/a.wav", "--model", "solaria-3"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}
	if postedBody["model"] != "solaria-3" {
		t.Fatalf("model = %v", postedBody["model"])
	}
	if _, ok := postedBody["language_config"]; ok {
		t.Fatalf("language_config should be omitted, got %#v", postedBody["language_config"])
	}
}

func TestTranscribeCommand_spaceSeparatedLanguage(t *testing.T) {
	withTempHome(t)
	t.Setenv(envGladiaAPIKey, "k")

	cases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "after source",
			args: []string{"transcribe", "https://example.com/a.wav", "--language", "en", "fr"},
			want: "en,fr",
		},
		{
			name: "before source",
			args: []string{"transcribe", "--language", "en", "fr", "meeting.wav"},
			want: "en,fr",
		},
		{
			name: "quoted value",
			args: []string{"transcribe", "https://example.com/a.wav", "--language", "en fr"},
			want: "en,fr",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := newRootCmd()
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), "comma-separated") || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
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
		case r.Method == http.MethodPost && r.URL.Path == "/v2/pre-recorded":
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
		case r.Method == http.MethodPost && r.URL.Path == "/v2/pre-recorded":
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

	run := func(args ...string) {
		t.Helper()
		postedBody = nil
		cmd := newRootCmd()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs(append([]string{"transcribe", "https://example.com/audio.wav"}, args...))
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	}

	assertCodeSwitchingPayload := func(t *testing.T) {
		t.Helper()
		lc, ok := postedBody["language_config"].(map[string]interface{})
		if !ok {
			t.Fatalf("language_config missing in %#v", postedBody)
		}
		if lc["code_switching"] != true {
			t.Fatalf("code_switching = %v, want true", lc["code_switching"])
		}
		langs, _ := lc["languages"].([]interface{})
		if len(langs) != 0 {
			t.Fatalf("languages = %v, want empty slice", lc["languages"])
		}
	}

	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "--code-switching", args: []string{"--code-switching"}},
		{name: "--cs", args: []string{"--cs"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			run(tc.args...)
			assertCodeSwitchingPayload(t)
		})
	}
}

func TestTranscribeCommand_modelRequestBody(t *testing.T) {
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
		case r.Method == http.MethodPost && r.URL.Path == "/v2/pre-recorded":
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

	run := func(args ...string) {
		t.Helper()
		postedBody = nil
		cmd := newRootCmd()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs(append([]string{"transcribe", "https://example.com/audio.wav"}, args...))
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	}

	t.Run("without --model", func(t *testing.T) {
		run()
		if _, ok := postedBody["model"]; ok {
			t.Fatalf("model present: %#v", postedBody)
		}
	})

	t.Run("with --model solaria-3", func(t *testing.T) {
		run("--model", "solaria-3", "--language", "en")
		if postedBody["model"] != "solaria-3" {
			t.Fatalf("model = %v", postedBody["model"])
		}
	})
}

func TestTranscribeCommand_diarizationRequestBody(t *testing.T) {
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
		case r.Method == http.MethodPost && r.URL.Path == "/v2/pre-recorded":
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

	run := func(args ...string) {
		t.Helper()
		postedBody = nil
		cmd := newRootCmd()
		cmd.SetOut(io.Discard)
		cmd.SetErr(io.Discard)
		cmd.SetArgs(append([]string{"transcribe", "https://example.com/audio.wav"}, args...))
		if err := cmd.Execute(); err != nil {
			t.Fatalf("execute: %v", err)
		}
	}

	t.Run("without --diarize", func(t *testing.T) {
		run()
		if _, ok := postedBody["diarization"]; ok {
			t.Fatalf("diarization present: %#v", postedBody)
		}
		if _, ok := postedBody["diarization_config"]; ok {
			t.Fatalf("diarization_config present: %#v", postedBody)
		}
	})

	t.Run("with --diarize", func(t *testing.T) {
		run("--diarize")
		if postedBody["diarization"] != true {
			t.Fatalf("diarization = %v", postedBody["diarization"])
		}
		cfg, ok := postedBody["diarization_config"].(map[string]interface{})
		if !ok {
			t.Fatalf("diarization_config missing: %#v", postedBody)
		}
		if cfg["min_speakers"] != float64(1) || cfg["max_speakers"] != float64(8) {
			t.Fatalf("speaker range = %v", cfg)
		}
		if _, ok := cfg["number_of_speakers"]; ok {
			t.Fatalf("number_of_speakers should be omitted: %v", cfg)
		}
	})
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
		case r.Method == http.MethodPost && r.URL.Path == "/v2/pre-recorded":
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
