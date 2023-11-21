package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	GLADIA_AUDIO_API_URL = "https://api.gladia.io/audio/text/audio-transcription/"
	GLADIA_VIDEO_API_URL = "https://api.gladia.io/video/text/video-transcription/"
	CONFIG_FILENAME      = ".gladia"
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

var LanguageList = []string{
	"afrikaans",
	"albanian",
	"amharic",
	"arabic",
	"armenian",
	"assamese",
	"azerbaijani",
	"bashkir",
	"basque",
	"belarusian",
	"bengali",
	"bosnian",
	"breton",
	"bulgarian",
	"catalan",
	"chinese",
	"croatian",
	"czech",
	"danish",
	"dutch",
	"english",
	"estonian",
	"faroese",
	"finnish",
	"french",
	"galician",
	"georgian",
	"german",
	"greek",
	"gujarati",
	"haitian creole",
	"hausa",
	"hawaiian",
	"hebrew",
	"hindi",
	"hungarian",
	"icelandic",
	"indonesian",
	"italian",
	"japanese",
	"javanese",
	"kannada",
	"kazakh",
	"khmer",
	"korean",
	"lao",
	"latin",
	"latvian",
	"lingala",
	"lithuanian",
	"luxembourgish",
	"macedonian",
	"malagasy",
	"malay",
	"malayalam",
	"maltese",
	"maori",
	"marathi",
	"mongolian",
	"myanmar",
	"nepali",
	"norwegian",
	"nynorsk",
	"occitan",
	"pashto",
	"persian",
	"polish",
	"portuguese",
	"punjabi",
	"romanian",
	"russian",
	"sanskrit",
	"serbian",
	"shona",
	"sindhi",
	"sinhala",
	"slovak",
	"slovenian",
	"somali",
	"spanish",
	"sundanese",
	"swahili",
	"swedish",
	"tagalog",
	"tajik",
	"tamil",
	"tatar",
	"telugu",
	"thai",
	"tibetan",
	"turkish",
	"turkmen",
	"ukrainian",
	"urdu",
	"uzbek",
	"vietnamese",
	"welsh",
	"wolof",
	"xhosa",
	"yiddish",
	"yoruba",
}

var TranslationList = []string{
	"afrikaans",
	"albanian",
	"amharic",
	"arabic",
	"armenian",
	"asturian",
	"azerbaijani",
	"bashkir",
	"belarusian",
	"bengali",
	"bosnian",
	"breton",
	"bulgarian",
	"burmese",
	"catalan",
	"cebuano",
	"chinese",
	"croatian",
	"czech",
	"danish",
	"dutch",
	"english",
	"estonian",
	"finnish",
	"flemish",
	"french",
	"western frisian",
	"fulah",
	"gaelic",
	"galician",
	"ganda",
	"georgian",
	"german",
	"greek",
	"gujarati",
	"haitian",
	"haitian creole",
	"hausa",
	"hebrew",
	"hindi",
	"hungarian",
	"icelandic",
	"igbo",
	"iloko",
	"indonesian",
	"irish",
	"italian",
	"japanese",
	"javanese",
	"kannada",
	"kazakh",
	"khmer",
	"korean",
	"lao",
	"latvian",
	"lingala",
	"lithuanian",
	"luxembourgish",
	"macedonian",
	"malagasy",
	"malay",
	"malayalam",
	"marathi",
	"moldavian",
	"moldovan",
	"mongolian",
	"nepali",
	"norwegian",
	"occitan",
	"oriya",
	"panjabi",
	"pashto",
	"persian",
	"polish",
	"portuguese",
	"pushto",
	"romanian",
	"russian",
	"serbian",
	"sindhi",
	"sinhala",
	"slovak",
	"slovenian",
	"somali",
	"spanish",
	"sundanese",
	"swahili",
	"swati",
	"swedish",
	"tagalog",
	"tamil",
	"thai",
	"tswana",
	"turkish",
	"ukrainian",
	"urdu",
	"uzbek",
	"valencian",
	"vietnamese",
	"welsh",
	"wolof",
	"xhosa",
	"yiddish",
	"yoruba",
}

type TranscriptionOptions struct {
	AudioURL                string `json:"audio_url"`
	AudioFile               string `json:"audio_file"`
	LanguageBehaviour       string `json:"language_behaviour"`
	Language                string `json:"language"`
	TranscriptionHint       string `json:"transcription_hint"`
	NoiseReduction          bool   `json:"noise_reduction"`
	Diarization             bool   `json:"diarization"`
	DiarizationMaxSpeakers  int    `json:"diarization_max_speakers"`
	DirectTranslate         bool   `json:"direct_translate"`
	DirectTranslateLanguage string `json:"direct_translate_language"`
	Summarization           bool   `json:"summarization"`
	OutputFormat            string `json:"output_format"`
	IsVideo                 string `json:"is_video"`
	LanguageList            bool   `json:"-"`
	TranslationList         bool   `json:"-"`
	GladiaKey               string `json:"-"`
	SaveGladiaKey           bool   `json:"-"`
}

type Prediction struct {
	Words         []Word      `json:"words"`
	Transcription string      `json:"transcription"`
	Language      string      `json:"language"`
	TimeBegin     float64     `json:"time_begin"`
	TimeEnd       float64     `json:"time_end"`
	Speaker       interface{} `json:"speaker"`
	Channel       string      `json:"channel"`
}

type Word struct {
	Word       string  `json:"word"`
	TimeBegin  float64 `json:"time_begin"`
	TimeEnd    float64 `json:"time_end"`
	Confidence float64 `json:"confidence"`
}

type ApiResponse struct {
	Prediction    []Prediction `json:"prediction"`
	PredictionRaw struct {
		Summarization string `json:"summarization"`
	} `json:"prediction_raw"`
}

type OtherResponse struct {
	Prediction    string `json:"prediction"`
	PredictionRaw struct {
		Summarization string `json:"summarization"`
	} `json:"prediction_raw"`
}

func getMIMEType(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil {
		return "", err
	}

	fileType := http.DetectContentType(buffer)

	return fileType, nil
}

func saveGladiaKeyToFile(gladiaKey string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, CONFIG_FILENAME)

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(gladiaKey + "\n")
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	fmt.Printf("Gladia API key saved to %s\n", configPath)
	return nil
}

func getGladiaKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, CONFIG_FILENAME)

	file, err := os.Open(configPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	gladiaKey, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	gladiaKey = strings.TrimSpace(gladiaKey)
	return gladiaKey, nil
}

func transcribe(options TranscriptionOptions) error {
	if options.GladiaKey == "" {
		gladiaKey, err := getGladiaKey()
		if err != nil {
			fmt.Println("Error: Gladia API key not found.")
			return err
		}
		options.GladiaKey = gladiaKey
	}

	if options.SaveGladiaKey {
		err := saveGladiaKeyToFile(options.GladiaKey)
		if err != nil {
			return err
		}
	}

	if options.GladiaKey == "" && !options.SaveGladiaKey {
		fmt.Println("Error: Gladia API key not found.")
		fmt.Println("Please provide your Gladia API key using --gladia-key or save it using --save-gladia-key.")
		return nil
	}

	if options.LanguageList {
		fmt.Println("Available Languages for Transcription:")
		for _, language := range LanguageList {
			fmt.Println(language)
		}
		return nil
	}

	if options.TranslationList {
		fmt.Println("Available Languages for Translation:")
		for _, language := range TranslationList {
			fmt.Println(language)
		}
		return nil
	}

	if !options.SaveGladiaKey {
		if options.GladiaKey != "" {
			if options.DirectTranslate && options.DirectTranslateLanguage == "" {
				fmt.Println("Error: --direct-translate-language is required when using --direct-translate.")
				return nil
			}

			if options.AudioURL == "" && options.AudioFile == "" {
				fmt.Println("Error: --audio-url or --audio-file is required.")
				return nil
			}

			client := &http.Client{}
			bodyWriter := &bytes.Buffer{}
			writer := multipart.NewWriter(bodyWriter)

			var urlField string = "audio_url"
			if options.IsVideo == "true" {
				urlField = "video_url"
			}

			if options.AudioURL != "" {
				err := addURLField(writer, urlField, options.AudioURL)
				if err != nil {
					fmt.Println("Error adding URL field:", err)
					return err
				}
			}

			if options.AudioFile != "" {
				err := addFileField(writer, options.AudioFile)
				if err != nil {
					fmt.Println("Error adding file field:", err)
					return err
				}
			}

			err := addStringField(writer, "language_behaviour", options.LanguageBehaviour)
			if err != nil {
				fmt.Println("Error adding language_behaviour field:", err)
				return err
			}

			err = addStringField(writer, "language", options.Language)
			if err != nil {
				fmt.Println("Error adding language field:", err)
				return err
			}

			err = addStringField(writer, "transcription_hint", options.TranscriptionHint)
			if err != nil {
				fmt.Println("Error adding transcription_hint field:", err)
				return err
			}

			err = addBoolField(writer, "noise_reduction", options.NoiseReduction)
			if err != nil {
				fmt.Println("Error adding noise_reduction field:", err)
				return err
			}

			err = addBoolField(writer, "toggle_diarization", options.Diarization)
			if err != nil {
				fmt.Println("Error adding toggle_diarization field:", err)
				return err
			}

			err = addIntField(writer, "diarization_max_speakers", options.DiarizationMaxSpeakers)
			if err != nil {
				fmt.Println("Error adding diarization_max_speakers field:", err)
				return err
			}

			err = addBoolField(writer, "toggle_direct_translate", options.DirectTranslate)
			if err != nil {
				fmt.Println("Error adding toggle_direct_translate field:", err)
				return err
			}

			err = addStringField(writer, "target_translation_language", options.DirectTranslateLanguage)
			if err != nil {
				fmt.Println("Error adding target_translation_language field:", err)
				return err
			}

			err = addBoolField(writer, "toggle_summarization", options.Summarization)
			if err != nil {
				fmt.Println("Error adding toggle_summarization field:", err)
				return err
			}

			if options.OutputFormat == "table" || options.OutputFormat == "csv" {
				err = addStringField(writer, "output_format", "json")
				if err != nil {
					fmt.Println("Error adding output_format field:", err)
					return err
				}
			} else {
				err = addStringField(writer, "output_format", options.OutputFormat)
				if err != nil {
					fmt.Println("Error adding output_format field:", err)
					return err
				}
			}

			err = writer.Close()
			if err != nil {
				fmt.Println("Error closing multipart writer:", err)
				return err
			}

			contentType := writer.FormDataContentType()

			// is IsVideo use the GLADIA_VIDEO_API_URL
			var GLADIA_API_URL string
			if options.IsVideo == "true" {
				GLADIA_API_URL = GLADIA_VIDEO_API_URL
			} else {
				GLADIA_API_URL = GLADIA_AUDIO_API_URL
			}
			req, err := http.NewRequest("POST", GLADIA_API_URL, bodyWriter)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return err
			}

			req.Header.Set("Accept", "application/json")
			req.Header.Set("X-Gladia-Key", options.GladiaKey)
			req.Header.Set("Content-Type", contentType)

			// Read the request body
			bodyBytes, err := ioutil.ReadAll(req.Body)
			if err != nil {
				fmt.Println("Error reading request body:", err)
				return err
			}

			// Restore the request body for sending the request
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error: %d - %s\n", resp.StatusCode, resp.Status)
				return nil
			}

			if options.OutputFormat == "table" {
				var apiResponse ApiResponse
				err = json.NewDecoder(resp.Body).Decode(&apiResponse)
				if err != nil {
					fmt.Println("Error decoding response:", err)
					return err
				}
				printTable(apiResponse.Prediction, options)
				if options.Summarization {
					fmt.Println()
					fmt.Println("=======")
					fmt.Println("Summary")
					fmt.Println("=======")
					fmt.Println(apiResponse.PredictionRaw.Summarization)
				}
			} else if options.OutputFormat == "csv" {
				var apiResponse ApiResponse
				err = json.NewDecoder(resp.Body).Decode(&apiResponse)
				if err != nil {
					fmt.Println("Error decoding response:", err)
					return err
				}
				printCSV(apiResponse.Prediction, options)
				if options.Summarization {
					fmt.Println()
					fmt.Println("=======")
					fmt.Println("Summary")
					fmt.Println("=======")
					fmt.Println(apiResponse.PredictionRaw.Summarization)
				}
			} else if options.OutputFormat == "srt" || options.OutputFormat == "vtt" || options.OutputFormat == "txt" {
				// only print the transcription
				var otherResponse OtherResponse
				err = json.NewDecoder(resp.Body).Decode(&otherResponse)
				if err != nil {
					fmt.Println("Error decoding response:", err)
					return err
				}
				fmt.Print(otherResponse.Prediction)
			} else {
				var apiResponse ApiResponse
				err = json.NewDecoder(resp.Body).Decode(&apiResponse)
				if err != nil {
					fmt.Println("Error decoding response:", err)
					return err
				}
				jsonBytes, err := json.MarshalIndent(apiResponse, "", "  ")
				if err != nil {
					fmt.Println("Error encoding JSON:", err)
					return err
				}
				fmt.Println(string(jsonBytes))
			}
		} else {
			fmt.Println("Error: Gladia API key not found.")
		}
	}

	return nil
}

func addURLField(writer *multipart.Writer, fieldName string, url string) error {
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, fieldName))
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(url))
	if err != nil {
		return err
	}

	return nil
}

func addFileField(writer *multipart.Writer, file string) error {
	fileType, err := getMIMEType(file)
	if err != nil {
		return err
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "audio", filepath.Base(file)))
	partHeader.Set("Content-Type", fileType)
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, f)
	if err != nil {
		return err
	}

	return nil
}

func addStringField(writer *multipart.Writer, fieldName string, value string) error {
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, fieldName))
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(value))
	if err != nil {
		return err
	}

	return nil
}

func addBoolField(writer *multipart.Writer, fieldName string, value bool) error {
	var stringValue string
	if value {
		stringValue = "true"
	} else {
		stringValue = "false"
	}

	return addStringField(writer, fieldName, stringValue)
}

func addIntField(writer *multipart.Writer, fieldName string, value int) error {
	return addStringField(writer, fieldName, fmt.Sprintf("%d", value))
}

func printTable(predictions []Prediction, options TranscriptionOptions) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time Begin", "Time End", "Language", "Speaker", "Transcription"})

	for _, prediction := range predictions {
		transcription := prediction.Transcription
		if options.DirectTranslate && prediction.Transcription != "" {
			transcription = translateText(prediction.Transcription, options.DirectTranslateLanguage)
		}

		speakerStr := ""
		switch speaker := prediction.Speaker.(type) {
		case int:
			speakerStr = fmt.Sprintf("%d", speaker)
		case string:
			speakerStr = speaker
		default:
			// Handle the case where the speaker is neither an int nor a string
			// You can choose to set a default value or handle the error appropriately
		}

		table.Append([]string{
			fmt.Sprintf("%.2f", prediction.TimeBegin),
			fmt.Sprintf("%.2f", prediction.TimeEnd),
			prediction.Language,
			speakerStr,
			transcription,
		})
	}

	table.Render()
}

func printCSV(predictions []Prediction, options TranscriptionOptions) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	header := []string{"Time Begin", "Time End", "Language", "Speaker", "Transcription"}
	writer.Write(header)

	for _, prediction := range predictions {
		transcription := prediction.Transcription
		if options.DirectTranslate && prediction.Transcription != "" {
			transcription = translateText(prediction.Transcription, options.DirectTranslateLanguage)
		}
		speakerStr := ""
		switch speaker := prediction.Speaker.(type) {
		case int:
			speakerStr = fmt.Sprintf("%d", speaker)
		case string:
			speakerStr = speaker
		default:
			// Handle the case where the speaker is neither an int nor a string
			// You can choose to set a default value or handle the error appropriately
		}

		row := []string{
			fmt.Sprintf("%.2f", prediction.TimeBegin),
			fmt.Sprintf("%.2f", prediction.TimeEnd),
			prediction.Language,
			speakerStr,
			transcription,
		}

		writer.Write(row)
	}

	writer.Flush()
}

func translateText(text string, language string) string {
	// Replace this with your own translation code
	// translation := text + " (Translated to " + language + ")"
	translation := text

	return translation
}

func main() {
	audioURLPtr := flag.String("audio-url", "", "URL of the audio file")
	audioFilePtr := flag.String("audio-file", "", "Path to the audio file")
	languageBehaviourPtr := flag.String("language-behaviour", "automatic multiple languages", "Language behavior (manual, automatic single language, automatic multiple languages)")
	languagePtr := flag.String("language", "english", "Language for transcription")
	transcriptionHintPtr := flag.String("transcription-hint", "", "Transcription hint")
	noiseReductionPtr := flag.Bool("noise-reduction", false, "Enable noise reduction")
	diarizationPtr := flag.Bool("diarization", false, "Enable diarization")
	diarizationMaxSpeakersPtr := flag.Int("diarization-max-speakers", 0, "Maximum number of speakers for diarization")
	directTranslatePtr := flag.Bool("direct-translate", false, "Enable direct translation")
	directTranslateLanguagePtr := flag.String("direct-translate-language", "", "Language for direct translation")
	summarizationPtr := flag.Bool("summarization", false, "Enable summarization")
	outputFormatPtr := flag.String("output-format", "table", "Output format (table, csv, json, srt, vtt, txt)")
	languageListPtr := flag.Bool("transcription-language-list", false, "List available languages for transcription")
	translationListPtr := flag.Bool("translation-language-list", false, "List available languages for translation")
	gladiaKeyPtr := flag.String("gladia-key", "", "Gladia API key")
	saveGladiaKeyPtr := flag.Bool("save-gladia-key", false, "Save Gladia API key")

	flag.Parse()

	options := TranscriptionOptions{
		AudioURL:                *audioURLPtr,
		AudioFile:               *audioFilePtr,
		LanguageBehaviour:       *languageBehaviourPtr,
		Language:                *languagePtr,
		TranscriptionHint:       *transcriptionHintPtr,
		NoiseReduction:          *noiseReductionPtr,
		Diarization:             *diarizationPtr,
		DiarizationMaxSpeakers:  *diarizationMaxSpeakersPtr,
		DirectTranslate:         *directTranslatePtr,
		DirectTranslateLanguage: *directTranslateLanguagePtr,
		Summarization:           *summarizationPtr,
		OutputFormat:            *outputFormatPtr,
		LanguageList:            *languageListPtr,
		TranslationList:         *translationListPtr,
		GladiaKey:               *gladiaKeyPtr,
		SaveGladiaKey:           *saveGladiaKeyPtr,
	}

	err := transcribe(options)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
