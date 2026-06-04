package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadFile_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/upload/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("x-gladia-key") != "secret" {
			t.Fatalf("api key header = %q", r.Header.Get("x-gladia-key"))
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"audio_url": "https://cdn/audio.wav"})
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.wav")
	if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := NewGladiaClient("secret", false)
	c.GladiaEndpoint = server.URL

	url, err := c.UploadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://cdn/audio.wav" {
		t.Fatalf("got %q", url)
	}
}

func TestUploadFile_apiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.wav")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	c := NewGladiaClient("k", false)
	c.GladiaEndpoint = server.URL
	if _, err := c.UploadFile(path); err == nil {
		t.Fatal("expected error")
	}
}

func TestTranscribeAudioURL_success(t *testing.T) {
	done := TranscriptionResult{Status: "done"}
	done.Result.Transcription.FullTranscript = "done text"
	doneJSON, _ := json.Marshal(done)

	var posted TranscriptionRequest
	server := httptest.NewServer(nil)
	base := server.URL
	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/transcription/":
			_ = json.NewDecoder(r.Body).Decode(&posted)
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(TranscriptionResponse{ResultURL: base + "/poll"})
		case r.Method == http.MethodGet:
			_, _ = w.Write(doneJSON)
		default:
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
	})
	defer server.Close()

	c := NewGladiaClient("k", false)
	c.GladiaEndpoint = server.URL

	req := TranscriptionRequest{
		LanguageConfig: &LanguageConfig{
			Languages:     []string{"en", "fr"},
			CodeSwitching: true,
		},
	}
	result, err := c.TranscribeAudioURL("https://audio.example/x.wav", req)
	if err != nil {
		t.Fatal(err)
	}
	if result.Result.Transcription.FullTranscript != "done text" {
		t.Fatalf("got %q", result.Result.Transcription.FullTranscript)
	}
	if posted.AudioURL != "https://audio.example/x.wav" {
		t.Fatalf("posted audio_url = %q", posted.AudioURL)
	}
	if posted.LanguageConfig == nil || !posted.LanguageConfig.CodeSwitching {
		t.Fatalf("language_config = %+v", posted.LanguageConfig)
	}
}

func TestTranscribeAudioURL_apiValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message":            "validation failed",
			"validation_errors":  []string{"bad field"},
		})
	}))
	defer server.Close()

	c := NewGladiaClient("k", false)
	c.GladiaEndpoint = server.URL
	_, err := c.TranscribeAudioURL("https://a", TranscriptionRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "validation failed") {
		t.Fatalf("got %v", err)
	}
}

func TestPollForTranscriptionResult_errorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"error","result":{"transcription":{"full_transcript":"boom"}}}`))
	}))
	defer server.Close()

	c := NewGladiaClient("k", false)
	c.GladiaEndpoint = server.URL
	_, err := c.pollForTranscriptionResult(server.URL + "/status")
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("got %v", err)
	}
}

func TestDecodeAPIError(t *testing.T) {
	c := NewGladiaClient("k", false)
	rec := httptest.NewRecorder()
	rec.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(rec).Encode(map[string]interface{}{
		"message":           "bad request",
		"validation_errors": []string{"a", "b"},
	})
	err := c.decodeAPIError(rec.Result())
	if err == nil || !strings.Contains(err.Error(), "bad request") || !strings.Contains(err.Error(), "a; b") {
		t.Fatalf("got %v", err)
	}
}

func TestTranscriptionRequest_marshalDiarization(t *testing.T) {
	t.Run("omits diarization when disabled", func(t *testing.T) {
		data, err := json.Marshal(TranscriptionRequest{AudioURL: "https://a"})
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		if _, ok := body["diarization"]; ok {
			t.Fatalf("diarization present: %v", body)
		}
		if _, ok := body["diarization_config"]; ok {
			t.Fatalf("diarization_config present: %v", body)
		}
	})

	t.Run("includes range without number_of_speakers", func(t *testing.T) {
		data, err := json.Marshal(TranscriptionRequest{
			AudioURL:    "https://a",
			Diarization: true,
			DiarizationConfig: &DiarizationConfig{
				MinSpeakers: 1,
				MaxSpeakers: 8,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]interface{}
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		if body["diarization"] != true {
			t.Fatalf("diarization = %v", body["diarization"])
		}
		cfg, ok := body["diarization_config"].(map[string]interface{})
		if !ok {
			t.Fatalf("diarization_config = %v", body["diarization_config"])
		}
		if cfg["min_speakers"] != float64(1) || cfg["max_speakers"] != float64(8) {
			t.Fatalf("speaker range = %v", cfg)
		}
		if _, ok := cfg["number_of_speakers"]; ok {
			t.Fatalf("number_of_speakers should be omitted: %v", cfg)
		}
	})
}

func TestSummarizationConfig_validate(t *testing.T) {
	if err := (&SummarizationConfig{Type: "general"}).ValidateSummarizationType(); err != nil {
		t.Fatal(err)
	}
	if err := (&SummarizationConfig{Type: "nope"}).ValidateSummarizationType(); err == nil {
		t.Fatal("expected error")
	}
}
