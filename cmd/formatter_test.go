package main

import (
	"encoding/json"
	"strings"
	"testing"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
)

func TestPrintTXTTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintTXTTranscription(sampleTranscriptionResult())
	})
	if strings.TrimSpace(out) != "hello world" {
		t.Fatalf("got %q", out)
	}
}

func TestPrintTXTDiarizedTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintTXTDiarizedTranscription(sampleTranscriptionResult())
	})
	if !strings.Contains(out, "Speaker 0: hello world") {
		t.Fatalf("got %q", out)
	}
	if !strings.Contains(out, "Speaker 1: second line") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintJSONSimplifiedTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintJSONSimplifiedTranscription(sampleTranscriptionResult())
	})
	var items []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d", len(items))
	}
	if items[0]["transcription"] != "hello world" {
		t.Fatalf("first item = %v", items[0])
	}
}

func TestPrintJSONSimplifiedTranscription_empty(t *testing.T) {
	var empty gladia.TranscriptionResult
	out := captureStdout(t, func() {
		PrintJSONSimplifiedTranscription(empty)
	})
	if !strings.Contains(out, "No transcriptions available") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintJSONTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintJSONTranscription(sampleTranscriptionResult())
	})
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatal(err)
	}
	if _, ok := parsed["result"]; !ok {
		t.Fatalf("missing result key in %v", parsed)
	}
}

func TestPrintSRTTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintSRTTranscription(sampleTranscriptionResult())
	})
	if !strings.Contains(out, "1\n00:00:00,500 --> 00:00:02,250\nhello world") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintSRTDiarizedTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintSRTDiarizedTranscription(sampleTranscriptionResult())
	})
	if !strings.Contains(out, "Speaker 0: hello world") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintVTTTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintVTTTranscription(sampleTranscriptionResult())
	})
	if !strings.HasPrefix(out, "WEBVTT\n") {
		t.Fatalf("got %q", out)
	}
	if !strings.Contains(out, "00:00:00.500 --> 00:00:02.250") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintVTTDiarizedTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintVTTDiarizedTranscription(sampleTranscriptionResult())
	})
	if !strings.Contains(out, "Speaker 1: second line") {
		t.Fatalf("got %q", out)
	}
}

func TestSecondsToSRTTimeFormat(t *testing.T) {
	tests := []struct {
		sec  float64
		want string
	}{
		{0.5, "00:00:00,500"},
		{65.25, "00:01:05,250"},
		{3661.001, "01:01:01,001"},
	}
	for _, tc := range tests {
		if got := secondsToSRTTimeFormat(tc.sec); got != tc.want {
			t.Errorf("secondsToSRTTimeFormat(%v) = %q, want %q", tc.sec, got, tc.want)
		}
	}
}

func TestSecondsToVTTTimeFormat(t *testing.T) {
	if got := secondsToVTTTimeFormat(65.25); got != "00:01:05.250" {
		t.Fatalf("got %q", got)
	}
}

func TestPrintCSVTranscription(t *testing.T) {
	out := captureStdout(t, func() {
		PrintCSVTranscription(sampleTranscriptionResult())
	})
	if !strings.HasPrefix(out, "time_begin, time_end, language, speaker, transcription\n") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintSummarization_withResult(t *testing.T) {
	result := sampleTranscriptionResult()
	summary := "short summary"
	result.Result.Summarization.Results = &summary
	out := captureStdout(t, func() {
		PrintSummarization(result)
	})
	if !strings.Contains(out, "short summary") {
		t.Fatalf("got %q", out)
	}
}

func TestPrintSummarization_empty(t *testing.T) {
	out := captureStdout(t, func() {
		PrintSummarization(sampleTranscriptionResult())
	})
	if !strings.Contains(out, "No summarization results available") {
		t.Fatalf("got %q", out)
	}
}
