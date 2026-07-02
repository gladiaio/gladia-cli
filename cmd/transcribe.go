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
  gladia transcribe call.wav -cs --language en -o json
  gladia transcribe call.wav --diarize -o srt
  gladia transcribe podcast.mp3 --model solaria-3 --language en
  gladia transcribe https://example.com/audio.mp3 -o json`,
		Args: validateTranscribeArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateOutputFormat(outputFormat); err != nil {
				return err
			}

			if err := validateModel(modelFlag); err != nil {
				return err
			}

			if err := validateLanguageFlag(languageFlag); err != nil {
				return err
			}

			langs, err := types.ParseLanguages(languageFlag)
			if err != nil {
				return err
			}

			codeSwitchSet := cmd.Flags().Changed("code-switching") || cmd.Flags().Changed("cs")
			if err := validateModelConfig(modelFlag, langs, codeSwitchSet, codeSwitching); err != nil {
				return err
			}
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
				Model:          normalizeModel(modelFlag),
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
	cmd.Flags().StringVar(&languageFlag, "language", "", "Expected language(s), comma-separated (e.g. en or en,fr,de); does not enable code switching")
	const codeSwitchingUsage = "Re-detect language on each utterance (for mixed-language audio; solaria-1 only)"
	cmd.Flags().BoolVar(&codeSwitching, "cs", false, codeSwitchingUsage)
	cmd.Flags().BoolVar(&codeSwitching, "code-switching", false, codeSwitchingUsage)
	cmd.Flags().Lookup("cs").Hidden = true
	cmd.Flags().Lookup("code-switching").Hidden = true
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show progress while transcribing")
	cmd.Flags().BoolVar(&diarization, "diarize", false, "Enable speaker diarization")
	cmd.Flags().StringVar(&modelFlag, "model", "", "STT model: solaria-1 or solaria-3 (solaria-3 accepts at most one --language: en, fr, de, es, or it)")

	cmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
      -cs, --code-switching   — ` + codeSwitchingUsage + `
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)

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
	model = normalizeModel(model)
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

var solaria3Languages = map[types.Language]bool{
	types.LanguageEn: true,
	types.LanguageFr: true,
	types.LanguageDe: true,
	types.LanguageEs: true,
	types.LanguageIt: true,
}

func validateModelConfig(model string, langs []types.Language, codeSwitchSet, codeSwitching bool) error {
	model = normalizeModel(model)
	if model != "solaria-3" {
		return nil
	}
	if codeSwitchSet && codeSwitching {
		return fmt.Errorf("solaria-3 does not support code switching (use solaria-1 instead)")
	}
	switch len(langs) {
	case 0:
		return nil
	case 1:
		if !solaria3Languages[langs[0]] {
			return fmt.Errorf("solaria-3 does not support language %q (use en, fr, de, es, or it)", langs[0])
		}
		return nil
	default:
		codes := make([]string, len(langs))
		for i, lang := range langs {
			codes[i] = string(lang)
		}
		return fmt.Errorf("solaria-3 accepts only one language, got %d (%s); use solaria-1 for multi-language", len(langs), strings.Join(codes, ", "))
	}
}

func normalizeModel(model string) string {
	model = strings.TrimSpace(strings.ToLower(model))
	return strings.ReplaceAll(model, " ", "-")
}

func validateTranscribeArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return nil
	}

	langFlag, _ := cmd.Flags().GetString("language")
	langFlag = strings.TrimSpace(langFlag)

	// gladia transcribe --language en fr meeting.wav
	if len(args) == 2 && isKnownLanguageCode(args[0]) && !isKnownLanguageCode(args[1]) && langFlag != "" {
		return spaceSeparatedLanguageError(joinLanguageCodes(langFlag, args[0]))
	}

	// gladia transcribe meeting.wav --language en fr
	var extraLangs []string
	for _, arg := range args[1:] {
		if isKnownLanguageCode(arg) {
			extraLangs = append(extraLangs, arg)
		}
	}
	if langFlag != "" && len(extraLangs) > 0 {
		return spaceSeparatedLanguageError(joinLanguageCodes(append([]string{langFlag}, extraLangs...)...))
	}

	return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
}

func validateLanguageFlag(s string) error {
	s = strings.TrimSpace(s)
	if s == "" || strings.Contains(s, ",") {
		return nil
	}
	if strings.Contains(s, " ") {
		parts := strings.Fields(s)
		if len(parts) > 1 && allKnownLanguageCodes(parts) {
			return spaceSeparatedLanguageError(parts)
		}
	}
	return nil
}

func spaceSeparatedLanguageError(codes []string) error {
	normalized := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.ToLower(strings.TrimSpace(code))
		if code != "" {
			normalized = append(normalized, code)
		}
	}
	return fmt.Errorf("--language expects comma-separated codes (e.g. --language %s), not spaces", strings.Join(normalized, ","))
}

func joinLanguageCodes(codes ...string) []string {
	out := make([]string, 0, len(codes))
	for _, code := range codes {
		code = strings.ToLower(strings.TrimSpace(code))
		if code != "" {
			out = append(out, code)
		}
	}
	return out
}

func allKnownLanguageCodes(codes []string) bool {
	for _, code := range codes {
		if !isKnownLanguageCode(code) {
			return false
		}
	}
	return len(codes) > 0
}

func isKnownLanguageCode(code string) bool {
	code = strings.TrimSpace(code)
	if code == "" {
		return false
	}
	_, err := types.ParseLanguage(code)
	return err == nil
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
