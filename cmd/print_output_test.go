package main

import (
	"bytes"
	"testing"
)

func TestPrintTranscriptionResult_formats(t *testing.T) {
	result := sampleTranscriptionResult()
	formats := []struct {
		format  string
		diarize bool
		check   func(string) bool
	}{
		{"text", false, func(s string) bool { return bytes.Contains([]byte(s), []byte("hello world")) }},
		{"txt", false, func(s string) bool { return bytes.Contains([]byte(s), []byte("hello world")) }},
		{"text", true, func(s string) bool { return bytes.Contains([]byte(s), []byte("Speaker 0:")) }},
		{"json", false, func(s string) bool { return bytes.Contains([]byte(s), []byte(`"transcription"`)) }},
		{"json-full", false, func(s string) bool { return bytes.Contains([]byte(s), []byte(`"full_transcript"`)) }},
		{"srt", false, func(s string) bool { return bytes.Contains([]byte(s), []byte("-->")) }},
		{"srt", true, func(s string) bool { return bytes.Contains([]byte(s), []byte("Speaker 0:")) }},
		{"vtt", false, func(s string) bool { return bytes.HasPrefix([]byte(s), []byte("WEBVTT")) }},
		{"vtt", true, func(s string) bool { return bytes.Contains([]byte(s), []byte("Speaker 1:")) }},
	}

	for _, tc := range formats {
		name := tc.format
		if tc.diarize {
			name += "+diarize"
		}
		t.Run(name, func(t *testing.T) {
			out := captureStdout(t, func() {
				printTranscriptionResult(result, tc.format, tc.diarize)
			})
			if !tc.check(out) {
				t.Fatalf("unexpected output: %q", out)
			}
		})
	}
}
