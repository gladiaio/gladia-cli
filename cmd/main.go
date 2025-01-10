package main

import (
	"flag"
	"fmt"
	"strings"

	gladia "github.com/gladiaio/gladia-cli/pkg/client"
	types "github.com/gladiaio/gladia-cli/pkg/client/types"
)

func main() {
	audioURLPtr := flag.String("audio-url", "", "URL of the audio file")
	audioFilePtr := flag.String("audio-file", "", "Path to the audio file")

	diarizationPtr := flag.Bool("diarization", false, "Enable diarization")
	diarizationMinSpeakersPtr := flag.Int("diarization-min-speakers", 1, "Minimum number of speakers")
	diarizationMaxSpeakersPtr := flag.Int("diarization-max-speakers", 8, "Maximum number of speakers")
	diarizationNumberOfSpeakersPtr := flag.Int("diarization-number-of-speakers", 4, "Number of speakers")

	enableCodeSwitchingPtr := flag.Bool("enable-code-switching", false, "Enable code switching")
	detectLanguagePtr := flag.Bool("detect-language", true, "Enable language detection")

	summarizationPtr := flag.Bool("summarization", false, "Enable summarization")
	summarizationTypePtr := flag.String("summarization-type", "general", "Summarization type")

	customVocabularyPtr := flag.String("custom-vocabulary", "", "Comma-separated list of custom vocabulary words")

	outputFormatPtr := flag.String("output-format", "table", "Output format (table, csv, json, etc.)")

	languageListPtr := flag.Bool("transcription-language-list", false, "List available languages")
	translationListPtr := flag.Bool("translation-language-list", false, "List translation languages")

	gladiaKeyPtr := flag.String("gladia-key", "", "Gladia API key")
	saveGladiaKeyPtr := flag.Bool("save-gladia-key", false, "Save Gladia API key")

	verbosePtr := flag.Bool("verbose", true, "Enable verbose printing (default=true)")

	flag.Parse()

	// 1) If we only intend to save the key (and no audio is passed), do so and skip the rest
	if *saveGladiaKeyPtr && *gladiaKeyPtr != "" && *audioURLPtr == "" && *audioFilePtr == "" {
		err := SaveGladiaKeyToFile(*gladiaKeyPtr)
		if err != nil {
			fmt.Printf("Error saving Gladia API key: %s\n", err)
		} else {
			fmt.Printf("Gladia API key saved successfully.\n")
		}
		// Immediately return so we don't prompt for audio
		return
	}

	// 2) Otherwise, if user also provided audio but wants to save the key, save it but continue
	if *saveGladiaKeyPtr && *gladiaKeyPtr != "" {
		err := SaveGladiaKeyToFile(*gladiaKeyPtr)
		if err != nil {
			fmt.Printf("Error saving Gladia API key: %s\n", err)
			return
		}
		fmt.Printf("Gladia API key saved successfully.\n")
	}

	// 3) If user did not provide --gladia-key, try reading from a stored file
	if *gladiaKeyPtr == "" {
		apiKey, err := GetGladiaKeyFromFile()
		if err != nil {
			fmt.Printf("Missing Gladia API key: %s\n", err)
			return
		}
		*gladiaKeyPtr = apiKey
	}

	client := gladia.NewGladiaClient(*gladiaKeyPtr, *verbosePtr)

	// 4) If just listing languages, do that and return
	if *languageListPtr {
		if _, err := types.DisplayAllInputLanguagesNames(); err != nil {
			fmt.Printf("Error getting languages: %s\n", err)
		}
		return
	}

	if *translationListPtr {
		if _, err := types.DisplayAllTargetLanguagesNames(); err != nil {
			fmt.Printf("Error getting languages: %s\n", err)
		}
		return
	}

	// 5) If no file or URL is provided, prompt for audio
	if *audioURLPtr == "" && *audioFilePtr == "" {
		fmt.Println("Please provide an audio URL or file path")
		return
	}

	// 6) Upload if there's a file, otherwise we'll use the provided URL
	var audioURL string
	var err error
	if *audioFilePtr != "" {
		audioURL, err = client.UploadFile(*audioFilePtr)
		if err != nil {
			fmt.Printf("Error uploading file: %s\n", err)
			return
		}
	} else {
		audioURL = *audioURLPtr
	}

	// 7) Build the transcription request
	transcriptionReq := gladia.TranscriptionRequest{
		Diarization: *diarizationPtr,
		DiarizationConfig: struct {
			MinSpeakers      int `json:"min_speakers"`
			MaxSpeakers      int `json:"max_speakers"`
			NumberOfSpeakers int `json:"number_of_speakers"`
		}{
			MinSpeakers:      *diarizationMinSpeakersPtr,
			MaxSpeakers:      *diarizationMaxSpeakersPtr,
			NumberOfSpeakers: *diarizationNumberOfSpeakersPtr,
		},
		EnableCodeSwitching: *enableCodeSwitchingPtr,
		DetectLanguage:      *detectLanguagePtr,
		Summarization:       *summarizationPtr,
		SummarizationConfig: &gladia.SummarizationConfig{
			Type: *summarizationTypePtr,
		},
		CustomVocabulary: strings.Split(*customVocabularyPtr, ","),
	}

	// 8) Transcribe and handle the result
	transcriptionResult, err := client.TranscribeAudioURL(audioURL, transcriptionReq)
	if err != nil {
		fmt.Printf("Transcription error: %s\n", err)
		return
	}

	if *verbosePtr {
		fmt.Println("Final transcription result:")
	}
	if *verbosePtr {
		fmt.Println(transcriptionResult.Result.Transcription.FullTranscript)
	} else {
		// If user doesn't want the final line, skip it or just print the bare transcript
		// e.g., fmt.Println(transcriptionResult.Result.Transcription.FullTranscript)
	}

	// 9) Format the output
	switch *outputFormatPtr {
	case "table":
		PrintTableTranscription(*transcriptionResult)
	case "csv":
		PrintCSVTranscription(*transcriptionResult)
	case "json":
		PrintJSONTranscription(*transcriptionResult)
	case "json-simplified":
		PrintJSONSimplifiedTranscription(*transcriptionResult)
	case "srt":
		PrintSRTTranscription(*transcriptionResult)
	case "srt-diarized":
		PrintSRTDiarizedTranscription(*transcriptionResult)
	case "vtt":
		PrintVTTTranscription(*transcriptionResult)
	case "vtt-diarized":
		PrintVTTDiarizedTranscription(*transcriptionResult)
	case "txt":
		PrintTXTTranscription(*transcriptionResult)
	case "txt-diarized":
		PrintTXTDiarizedTranscription(*transcriptionResult)
	case "summary":
		PrintSummarization(*transcriptionResult)
	default:
		PrintTableTranscription(*transcriptionResult)
	}
}
