package main

import (
	"flag"
	"fmt"
	"strings"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
	types "github.com/gladiaio/gladia-cli/pkg/client/types"
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
	var UploadResponse *gladia.UploadResponse
	flag.Parse()

	if *saveGladiaKeyPtr && *gladiaKeyPtr != "" {
		err := SaveGladiaKeyToFile(*gladiaKeyPtr)
		if err != nil {
			fmt.Printf("Error saving Gladia API key: %s\n", err)
			return
		}
	}

	if *gladiaKeyPtr == "" {
		apiKey, err := GetGladiaKeyFromFile()
		if err != nil {
			fmt.Printf("Missing Gladia API key: %s\n", err)
			return
		}

		*gladiaKeyPtr = apiKey
	}

	client := gladia.NewGladiaClient(*gladiaKeyPtr)

	if *languageListPtr {
		_, err := types.DisplayAllInputLanguagesNames()
		if err != nil {
			fmt.Printf("Error getting languages: %s\n", err)
			return
		}

		return
	}

	if *translationListPtr {
		_, err := types.DisplayAllTargetLanguagesNames()
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

		UploadResponse, err = client.UploadFile(*audioFilePtr)
		if err != nil {
			fmt.Printf("Error uploading file: %s\n", err)
			return
		}
	}

	var transcriptionRequest gladia.TranscriptionRequest
	transcriptionRequest.AudioURL = UploadResponse.AudioURL
	transcriptionRequest.Diarization = *diarizationPtr
	transcriptionRequest.DiarizationConfig.MinSpeakers = *diarizationMinSpeakersPtr
	transcriptionRequest.DiarizationConfig.MaxSpeakers = *diarizationMaxSpeakersPtr
	transcriptionRequest.DiarizationConfig.NumberOfSpeakers = *diarizationNumberOfSpeakersPtr
	transcriptionRequest.EnableCodeSwitching = *enableCodeSwitchingPtr
	transcriptionRequest.DetectLanguage = *detectLanguagePtr
	transcriptionRequest.Summarization = *summarizationPtr
	transcriptionRequest.SummarizationConfig = &gladia.SummarizationConfig{Type: *summarizationTypePtr}
	transcriptionRequest.CustomVocabulary = strings.Split(*customVocabularyPtr, ",")

	transcription, err := client.GetTranscription(transcriptionRequest)
	if err != nil {
		fmt.Printf("Error getting transcription: %s\n", err)
		return
	}
	println()
	switch *outputFormatPtr {
	case "table":
		PrintTableTranscription(*transcription)
	case "csv":
		PrintCSVTranscription(*transcription)
	case "json":
		PrintJSONTranscription(*transcription)
	case "json-simplified":
		PrintJSONSimplifiedTranscription(*transcription)
	case "srt":
		PrintSRTTranscription(*transcription)
	case "srt-diarized":
		PrintSRTDiarizedTranscription(*transcription)
	case "vtt":
		PrintVTTTranscription(*transcription)
	case "vtt-diarized":
		PrintVTTDiarizedTranscription(*transcription)
	case "txt":
		PrintTXTTranscription(*transcription)
	case "txt-diarized":
		PrintTXTDiarizedTranscription(*transcription)
	case "summary":
		PrintSummarization(*transcription)
	default:
		PrintTableTranscription(*transcription)
	}
}
