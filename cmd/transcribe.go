package main

import (
	"fmt"
	"os"
	"strings"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
	"github.com/spf13/cobra"
)

func newTranscribeCmd() *cobra.Command {
	var (
		outputFormat string
		verbose      bool
		diarization  bool
	)

	cmd := &cobra.Command{
		Use:   "transcribe [source]",
		Short: "Transcribe a local audio file or URL",
		Long: `Transcribe pre-recorded audio from a file path or http(s) URL.

Examples:
  gladia transcribe meeting.wav
  gladia transcribe audio.mp3 -o text
  gladia transcribe https://example.com/audio.mp3 -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := ResolveAPIKey(rootGladiaKey)
			if err != nil {
				return err
			}

			client := gladia.NewGladiaClient(key, verbose)

			audioURL, err := resolveAudioSource(client, args[0])
			if err != nil {
				return err
			}

			transcriptionReq := gladia.TranscriptionRequest{
				Diarization: diarization,
				DiarizationConfig: struct {
					MinSpeakers      int `json:"min_speakers"`
					MaxSpeakers      int `json:"max_speakers"`
					NumberOfSpeakers int `json:"number_of_speakers"`
				}{
					MinSpeakers:      1,
					MaxSpeakers:      8,
					NumberOfSpeakers: 4,
				},
				DetectLanguage: true,
			}

			result, err := client.TranscribeAudioURL(audioURL, transcriptionReq)
			if err != nil {
				return fmt.Errorf("transcription failed: %w", err)
			}

			printTranscriptionResult(*result, outputFormat)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, json-full")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show progress while transcribing")
	cmd.Flags().BoolVar(&diarization, "diarize", false, "Enable speaker diarization")

	return cmd
}

func isHTTPURL(s string) bool {
	lower := strings.ToLower(s)
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func resolveAudioSource(client *gladia.GladiaClient, source string) (string, error) {
	if isHTTPURL(source) {
		return source, nil
	}

	if _, err := os.Stat(source); err != nil {
		return "", fmt.Errorf("%q is not a URL and not a readable file: %w", source, err)
	}

	audioURL, err := client.UploadFile(source)
	if err != nil {
		return "", fmt.Errorf("upload file: %w", err)
	}
	return audioURL, nil
}

func printTranscriptionResult(result gladia.TranscriptionResult, format string) {
	switch format {
	case "text", "txt":
		PrintTXTTranscription(result)
	case "json":
		PrintJSONSimplifiedTranscription(result)
	case "json-full":
		PrintJSONTranscription(result)
	default:
		fmt.Fprintf(os.Stderr, "unknown output format %q (use text, json, or json-full)\n", format)
		os.Exit(1)
	}
}
