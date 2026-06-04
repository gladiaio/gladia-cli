package main

import (
	"os"
	"testing"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
)

func withTempHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	buf := make([]byte, 1<<20)
	n, _ := r.Read(buf)
	return string(buf[:n])
}

func sampleTranscriptionResult() gladia.TranscriptionResult {
	var result gladia.TranscriptionResult
	result.Result.Transcription.FullTranscript = "hello world"
	result.Result.Transcription.Utterances = gladia.Utterances{
		{
			Start:    0.5,
			End:      2.25,
			Language: "en",
			Speaker:  0,
			Text:     "hello world",
		},
		{
			Start:    3.0,
			End:      5.5,
			Language: "en",
			Speaker:  1,
			Text:     "second line",
		},
	}
	return result
}
