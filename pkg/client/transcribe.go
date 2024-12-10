package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// UploadResponse represents the JSON response structure from the API.
type UploadResponse struct {
	AudioURL      string `json:"audio_url"`
	AudioMetadata struct {
		ID            string  `json:"id"`
		Filename      string  `json:"filename"`
		Extension     string  `json:"extension"`
		Size          int     `json:"size"`
		AudioDuration float64 `json:"audio_duration"`
		NbChannels    int     `json:"nb_channels"`
	} `json:"audio_metadata"`
}

type TranscriptionRequest struct {
	AudioURL          string `json:"audio_url"`
	Diarization       bool   `json:"diarization"`
	DiarizationConfig struct {
		MinSpeakers      int `json:"min_speakers"`
		MaxSpeakers      int `json:"max_speakers"`
		NumberOfSpeakers int `json:"number_of_speakers"`
	} `json:"diarization_config"`
	EnableCodeSwitching bool                 `json:"enable_code_switching"`
	DetectLanguage      bool                 `json:"detect_language"`
	Summarization       bool                 `json:"summarization"`
	SummarizationConfig *SummarizationConfig `json:"summarization_config"`
	Translation         bool                 `json:"translation"`
	TranslationConfig   *TranslationConfig   `json:"translation_config"`
	CustomVocabulary    []string             `json:"custom_vocabulary"`
}

type TranslationConfig struct {
	Model           string   `json:"model"`
	TargetLanguages []string `json:"target_languages"`
}

// SummarizationConfig represents the configuration for summarization.
type SummarizationConfig struct {
	Type string `json:"type"`
}

// ValidateSummarizationType checks if the SummarizationConfig Type is valid.
func (sc *SummarizationConfig) ValidateSummarizationType() error {
	switch sc.Type {
	case "general", "bullet_points", "concise":
		return nil
	default:
		return fmt.Errorf("invalid summarization type: %s", sc.Type)
	}
}

type TranscriptionResponse struct {
	ID        string `json:"id"`
	ResultURL string `json:"result_url"`
}

type TranscriptionResult struct {
	ID            string `json:"id"`
	RequestID     string `json:"request_id"`
	Kind          string `json:"kind"`
	Status        string `json:"status"`
	RequestParams struct {
		AudioURL          string `json:"audio_url"`
		Translation       bool   `json:"translation"`
		TranslationConfig struct {
			Model           string   `json:"model"`
			TargetLanguages []string `json:"target_languages"`
		} `json:"translation_config"`
	} `json:"request_params"`
	Result struct {
		CustomPrompts struct {
			Error struct {
				Exception  string `json:"exception"`
				Message    string `json:"message"`
				StatusCode int    `json:"status_code"`
			} `json:"error"`
			ExecTime float64 `json:"exec_time"`
			IsEmpty  bool    `json:"is_empty"`
			Results  []struct {
				Error struct {
					Exception  string `json:"exception"`
					Message    string `json:"message"`
					StatusCode int    `json:"status_code"`
				} `json:"error"`
				ExecTime float64 `json:"exec_time"`
				IsEmpty  bool    `json:"is_empty"`
				Results  struct {
					Prompt   *string `json:"prompt"`
					Response *string `json:"response"`
				} `json:"results"`
				Success bool `json:"success"`
			} `json:"results"`
			Success bool `json:"success"`
		} `json:"custom_prompts"`
		Metadata struct {
			AudioDuration            float64 `json:"audio_duration"`
			NumberOfDistinctChannels int     `json:"number_of_distinct_channels"`
			BillingTime              float64 `json:"billing_time"`
			TranscriptionTime        float64 `json:"transcription_time"`
		} `json:"metadata"`
		Summarization struct {
			Error struct {
				Exception  string `json:"exception"`
				Message    string `json:"message"`
				StatusCode int    `json:"status_code"`
			} `json:"error"`
			ExecTime float64 `json:"exec_time"`
			IsEmpty  bool    `json:"is_empty"`
			Results  *string `json:"results"`
			Success  bool    `json:"success"`
		} `json:"summarization"`
		Transcription struct {
			FullTranscript string   `json:"full_transcript"`
			Languages      []string `json:"languages"`
			Subtitles      []struct {
				Format    string `json:"format"`
				Subtitles string `json:"subtitles"`
			} `json:"subtitles"`
			Utterances Utterances `json:"utterances"`
		} `json:"transcription"`
		Translation struct {
			Error struct {
				Exception  string `json:"exception"`
				Message    string `json:"message"`
				StatusCode int    `json:"status_code"`
			} `json:"error"`
			ExecTime float64 `json:"exec_time"`
			IsEmpty  bool    `json:"is_empty"`
			Results  []struct {
				FullTranscript string   `json:"full_transcript"`
				Languages      []string `json:"languages"`
				Subtitles      []struct {
					Format    string `json:"format"`
					Subtitles string `json:"subtitles"`
				} `json:"subtitles"`
				Utterances Utterances `json:"utterances"`
			} `json:"results"`
			Success bool `json:"success"`
		} `json:"translation"`
	} `json:"result"`
}

// Utterance represents a segment of transcribed text with metadata.
type Utterances []struct {
	Channel    int     `json:"channel"`
	Confidence float64 `json:"confidence"`
	End        float64 `json:"end"`
	Language   string  `json:"language"`
	Speaker    int     `json:"speaker"`
	Start      float64 `json:"start"`
	Text       string  `json:"text"`
	Words      []struct {
		Confidence float64 `json:"confidence"`
		End        float64 `json:"end"`
		Start      float64 `json:"start"`
		Word       string  `json:"word"`
	} `json:"words"`
}

// UploadFile uploads an audio file to the Gladia API and returns the upload response.
func (c *GladiaClient) UploadFile(filePath string) (*UploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Printf("Failed to close file %s: %v", filePath, cerr)
		}
	}()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio", filePath)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		log.Printf("Failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", c.GladiaEndpoint+"/v2/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-gladia-key", c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	var uploadResp UploadResponse
	if err = json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	return &uploadResp, nil
}

// Create and execute a new HTTP request with JSON body
func (c *GladiaClient) createAndExecuteRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-gladia-key", c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Poll for the transcription result
func (c *GladiaClient) pollForTranscriptionResult(resultURL string) (*TranscriptionResult, error) {
	spinner := []string{"-", "\\", "|", "/"}
	spinnerIndex := 0

	for {
		resp, err := c.createAndExecuteRequest("GET", resultURL, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result TranscriptionResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		if result.Status == "done" {
			fmt.Printf("\033[H\033[2J") // Clear the terminal screen
			fmt.Println("Transcription completed successfully.")
			fmt.Printf("\033[H\033[2J") // Clear the terminal screen
			return &result, nil
		}

		if result.Status == "error" {
			fmt.Printf("\033[H\033[2J") // Clear the terminal screen
			return nil, fmt.Errorf("transcription failed with error: %s", result.Result.Transcription.FullTranscript)
		}

		fmt.Printf("\rTranscription in progress... %s     (%s) (request_id: %s)", spinner[spinnerIndex], result.Status, result.RequestID)
		spinnerIndex = (spinnerIndex + 1) % len(spinner)

		time.Sleep(1 * time.Second)
	}
}

func (c *GladiaClient) GetTranscription(transcriptionRequest TranscriptionRequest) (*TranscriptionResult, error) {
	requestBody, err := json.Marshal(transcriptionRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.createAndExecuteRequest("POST", c.GladiaEndpoint+"/v2/transcription", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var respError struct {
			Message          string   `json:"message"`
			Path             string   `json:"path"`
			RequestID        string   `json:"request_id"`
			StatusCode       int      `json:"statusCode"`
			Timestamp        string   `json:"timestamp"`
			ValidationErrors []string `json:"validation_errors"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&respError); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %v", err)
		}

		errorMessage := fmt.Sprintf("Error message: %s \n Validation errors: %v", respError.Message, respError.ValidationErrors)
		println(errorMessage)
		return nil, fmt.Errorf("failed to request transcription, status code: %d %s", resp.StatusCode, respError.Message)
	}

	var transcriptionResponse TranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transcriptionResponse); err != nil {
		return nil, err
	}

	return c.pollForTranscriptionResult(transcriptionResponse.ResultURL)
}
