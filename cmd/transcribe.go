package main

import (
	"fmt"
	"os"
	"strings"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
	"github.com/gladiaio/gladia-cli/pkg/client/types"
	"github.com/spf13/cobra"
)

func newTranscribeCmd() *cobra.Command {
	var (
		outputFormat  string
		languageFlag  string
		codeSwitching bool
		verbose       bool
		diarization   bool
		modelFlag     string
	)

	cmd := &cobra.Command{
		Use:   "transcribe [source]",
		Short: "Transcribe a local audio file or URL",
		Long: `Transcribe pre-recorded audio from a file path or http(s) URL.

Examples:
  gladia transcribe meeting.wav
  gladia transcribe audio.mp3 -o text
  gladia transcribe podcast.mp3 --language en
  gladia transcribe interview.mp3 --code-switching
  gladia transcribe interview.mp3 --language en,fr,de
  gladia transcribe call.wav --code-switch --language en -o json
  gladia transcribe call.wav --diarize -o srt
  gladia transcribe podcast.mp3 --model solaria-3
  gladia transcribe https://example.com/audio.mp3 -o json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateOutputFormat(outputFormat); err != nil {
				return err
			}

			if err := validateModel(modelFlag); err != nil {
				return err
			}

			langs, err := types.ParseLanguages(languageFlag)
			if err != nil {
				return err
			}

			codeSwitchSet := cmd.Flags().Changed("code-switching") || cmd.Flags().Changed("code-switch")
			langConfig, err := buildLanguageConfig(langs, codeSwitching, codeSwitchSet)
			if err != nil {
				return err
			}

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
				Model:          modelFlag,
				LanguageConfig: langConfig,
				Diarization:    diarization,
			}
			if diarization {
				transcriptionReq.DiarizationConfig = &gladia.DiarizationConfig{
					MinSpeakers: 1,
					MaxSpeakers: 8,
				}
			}

			result, err := client.TranscribeAudioURL(audioURL, transcriptionReq)
			if err != nil {
				return fmt.Errorf("transcription failed: %w", err)
			}

			printTranscriptionResult(*result, outputFormat, diarization)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, json-full, srt, vtt")
	cmd.Flags().StringVar(&languageFlag, "language", "", "Optional ISO 639-1 code(s), comma-separated (e.g. en or en,fr,de)")
	cmd.Flags().BoolVar(&codeSwitching, "code-switching", false, "Enable code switching (detect language per utterance; independent of --language)")
	cmd.Flags().BoolVar(&codeSwitching, "code-switch", false, "Alias for --code-switching")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show progress while transcribing")
	cmd.Flags().BoolVar(&diarization, "diarize", false, "Enable speaker diarization")
	cmd.Flags().StringVar(&modelFlag, "model", "", "STT model: solaria-1 or solaria-3 (default: API default)")

	return cmd
}

func buildLanguageConfig(langs []types.Language, codeSwitching, codeSwitchSet bool) (*gladia.LanguageConfig, error) {
	if len(langs) == 0 && !codeSwitchSet {
		return nil, nil
	}

	codes := make([]string, len(langs))
	for i, lang := range langs {
		codes[i] = string(lang)
	}

	cfg := &gladia.LanguageConfig{Languages: codes}

	if codeSwitchSet {
		cfg.CodeSwitching = codeSwitching
	}

	return cfg, nil
}

func validateOutputFormat(format string) error {
	switch format {
	case "text", "txt", "json", "json-full", "srt", "vtt":
		return nil
	default:
		return fmt.Errorf("unknown output format %q (use text, json, json-full, srt, or vtt)", format)
	}
}

func validateModel(model string) error {
	if model == "" {
		return nil
	}
	switch model {
	case "solaria-1", "solaria-3":
		return nil
	default:
		return fmt.Errorf("unknown model %q (use solaria-1 or solaria-3)", model)
	}
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

func printTranscriptionResult(result gladia.TranscriptionResult, format string, diarize bool) {
	switch format {
	case "text", "txt":
		if diarize {
			PrintTXTDiarizedTranscription(result)
		} else {
			PrintTXTTranscription(result)
		}
	case "json":
		PrintJSONSimplifiedTranscription(result)
	case "json-full":
		PrintJSONTranscription(result)
	case "srt":
		if diarize {
			PrintSRTDiarizedTranscription(result)
		} else {
			PrintSRTTranscription(result)
		}
	case "vtt":
		if diarize {
			PrintVTTDiarizedTranscription(result)
		} else {
			PrintVTTTranscription(result)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown output format %q\n", format)
		os.Exit(1)
	}
}
