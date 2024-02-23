package main

import (
	"flag"
	"fmt"
	languages "gladia/api/languages"
	transcribe "gladia/api/transcribe"
	formatter "gladia/cli/formatter"
	key "gladia/cli/key"
	"strings"
)

const (
	GLADIA_API_URL  = "https://api.gladia.io/v2"
	CONFIG_FILENAME = ".gladia"
)

type Color struct {
	Reset     string
	Purple    string
	Cyan      string
	DarkCyan  string
	Blue      string
	Green     string
	Yellow    string
	Red       string
	Bold      string
	Underline string
}

var Colors = Color{
	Reset:     "\033[0m",
	Purple:    "\033[95m",
	Cyan:      "\033[96m",
	DarkCyan:  "\033[36m",
	Blue:      "\033[94m",
	Green:     "\033[92m",
	Yellow:    "\033[93m",
	Red:       "\033[91m",
	Bold:      "\033[1m",
	Underline: "\033[4m",
}

var Language languages.Language
var TargetLanguage languages.TargetLanguage
var UploadResponse transcribe.UploadResponse

func main() {
	audioURLPtr := flag.String("audio-url", "", "URL of the audio file")
	audioFilePtr := flag.String("audio-file", "", "Path to the audio file")
	// languageBehaviourPtr := flag.String("language-behaviour", "automatic multiple languages", "Language behavior (manual, automatic single language, automatic multiple languages)")
	// languagePtr := flag.String("language", "english", "Language for transcription")
	// transcriptionHintPtr := flag.String("transcription-hint", "", "Transcription hint")
	// noiseReductionPtr := flag.Bool("noise-reduction", false, "Enable noise reduction")
	diarizationPtr := flag.Bool("diarization", false, "Enable diarization")
	diarizationMinSpeakersPtr := flag.Int("diarization-min-speakers", 1, "Minimum number of speakers for diarization")
	diarizationMaxSpeakersPtr := flag.Int("diarization-max-speakers", 8, "Maximum number of speakers for diarization")
	diarizationNumberOfSpeakersPtr := flag.Int("diarization-number-of-speakers", 4, "Number of speakers for diarization")

	enableCodeSwitchingPtr := flag.Bool("enable-code-switching", false, "Enable code switching")
	detectLanguagePtr := flag.Bool("detect-language", true, "Enable language detection")

	summarizationPtr := flag.Bool("summarization", false, "Enable summarization")
	summarizationTypePtr := flag.String("summarization-type", "general", "Summarization type (general, bullet_points, concise)")

	customVocabularyPtr := flag.String("custom-vocabulary", "", "Custom vocabulary use a comma separated list of words")

	outputFormatPtr := flag.String("output-format", "table", "Output format (table, csv, json, json-simplified, srt, srt-diarized, vtt, vtt-diarized, txt, txt-diarized, summary)")

	languageListPtr := flag.Bool("transcription-language-list", false, "List available languages for transcription")
	translationListPtr := flag.Bool("translation-language-list", false, "List available languages for translation")

	gladiaKeyPtr := flag.String("gladia-key", "", "Gladia API key")
	saveGladiaKeyPtr := flag.Bool("save-gladia-key", false, "Save Gladia API key")
	var UploadResponse *transcribe.UploadResponse
	flag.Parse()

	if *saveGladiaKeyPtr && *gladiaKeyPtr != "" {
		err := key.SaveGladiaKeyToFile(*gladiaKeyPtr)
		if err != nil {
			fmt.Printf("Error saving Gladia API key: %s\n", err)
			return
		}
	}

	if *languageListPtr {
		_, err := languages.DisplayAllInputLanguagesNames()
		if err != nil {
			fmt.Printf("Error getting languages: %s\n", err)
			return
		}

		return
	}

	if *translationListPtr {
		_, err := languages.DisplayAllTargetLanguagesNames()
		if err != nil {
			fmt.Printf("Error getting languages: %s\n", err)
			return
		}

		return
	}

	if *audioURLPtr == "" && *audioFilePtr == "" {
		fmt.Println("Please provide an audio URL or file path")
		return
	}

	if *audioURLPtr != "" && *audioFilePtr != "" {
		fmt.Println("Please provide only an audio URL or file path")
		return
	}

	if *audioURLPtr != "" {
		UploadResponse.AudioURL = *audioURLPtr
	} else {
		if *audioFilePtr == "" {
			fmt.Println("Please provide an audio file path")
			return
		}
		var err error
		UploadResponse, err = transcribe.UploadFile(*audioFilePtr)
		if err != nil {
			fmt.Printf("Error uploading file: %s\n", err)
			return
		}
	}

	var transcriptionRequest transcribe.TranscriptionRequest
	transcriptionRequest.AudioURL = UploadResponse.AudioURL
	transcriptionRequest.Diarization = *diarizationPtr
	transcriptionRequest.DiarizationConfig.MinSpeakers = *diarizationMinSpeakersPtr
	transcriptionRequest.DiarizationConfig.MaxSpeakers = *diarizationMaxSpeakersPtr
	transcriptionRequest.DiarizationConfig.NumberOfSpeakers = *diarizationNumberOfSpeakersPtr
	transcriptionRequest.EnableCodeSwitching = *enableCodeSwitchingPtr
	transcriptionRequest.DetectLanguage = *detectLanguagePtr
	transcriptionRequest.Summarization = *summarizationPtr
	transcriptionRequest.SummarizationConfig = &transcribe.SummarizationConfig{Type: *summarizationTypePtr}
	transcriptionRequest.CustomVocabulary = strings.Split(*customVocabularyPtr, ",")

	transcription, err := transcribe.GetTranscription(transcriptionRequest)
	if err != nil {
		fmt.Printf("Error getting transcription: %s\n", err)
		return
	}
	println()
	switch *outputFormatPtr {
	case "table":
		formatter.PrintTableTranscription(*transcription)
	case "csv":
		formatter.PrintCSVTranscription(*transcription)
	case "json":
		formatter.PrintJSONTranscription(*transcription)
	case "json-simplified":
		formatter.PrintJSONSimplifiedTranscription(*transcription)
	case "srt":
		formatter.PrintSRTTranscription(*transcription)
	case "srt-diarized":
		formatter.PrintSRTDiarizedTranscription(*transcription)
	case "vtt":
		formatter.PrintVTTTranscription(*transcription)
	case "vtt-diarized":
		formatter.PrintVTTDiarizedTranscription(*transcription)
	case "txt":
		formatter.PrintTXTTranscription(*transcription)
	case "txt-diarized":
		formatter.PrintTXTDiarizedTranscription(*transcription)
	case "summary":
		formatter.PrintSummarization(*transcription)
	default:
		formatter.PrintTableTranscription(*transcription)
	}
}
