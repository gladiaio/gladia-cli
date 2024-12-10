package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TranscriptionItem represents a single transcription item in the list.
type TranscriptionItem struct {
	CreatedAt      string                 `json:"created_at"`
	CustomMetadata map[string]interface{} `json:"custom_metadata"`
	File           struct {
		AudioDuration    int    `json:"audio_duration"`
		Filename         string `json:"filename"`
		ID               string `json:"id"`
		NumberOfChannels int    `json:"number_of_channels"`
		Source           string `json:"source"`
	} `json:"file"`
	ID            string   `json:"id"`
	Kind          []string `json:"kind"`
	RequestID     string   `json:"request_id"`
	RequestParams struct {
		AudioURL            string `json:"audio_url"`
		CallbackURL         string `json:"callback_url"`
		ContextPrompt       string `json:"context_prompt"`
		CustomPrompts       bool   `json:"custom_prompts"`
		CustomPromptsConfig struct {
			Prompts []string `json:"prompts"`
		} `json:"custom_prompts_config"`
		CustomVocabulary  []string `json:"custom_vocabulary"`
		DetectLanguage    bool     `json:"detect_language"`
		Diarization       bool     `json:"diarization"`
		DiarizationConfig struct {
			MaxSpeakers      int `json:"max_speakers"`
			MinSpeakers      int `json:"min_speakers"`
			NumberOfSpeakers int `json:"number_of_speakers"`
		} `json:"diarization_config"`
		EnableCodeSwitching bool   `json:"enable_code_switching"`
		Language            string `json:"language"`
		Subtitles           bool   `json:"subtitles"`
		SubtitlesConfig     struct {
			Formats []string `json:"formats"`
		} `json:"subtitles_config"`
		Summarization       bool `json:"summarization"`
		SummarizationConfig struct {
			Type string `json:"type"`
		} `json:"summarization_config"`
		Translation       bool `json:"translation"`
		TranslationConfig struct {
			Model           string   `json:"model"`
			TargetLanguages []string `json:"target_languages"`
		} `json:"translation_config"`
	} `json:"request_params"`
	Status string `json:"status"`
}

// TranscriptionListResponse represents the response structure for listing transcriptions.
type TranscriptionListResponse struct {
	Current string              `json:"current"`
	First   string              `json:"first"`
	Items   []TranscriptionItem `json:"items"`
	Next    string              `json:"next"`
}

// ListTranscriptions fetches the list of transcriptions from the Gladia API.
func (c *GladiaClient) ListTranscriptions(offset, limit int, status, kind, date, beforeDate, afterDate string) (*TranscriptionListResponse, error) {
	req, err := http.NewRequest("GET", c.GladiaEndpoint+"/v2/transcription", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("offset", fmt.Sprintf("%d", offset))
	q.Add("limit", fmt.Sprintf("%d", limit))
	q.Add("status", status)
	q.Add("kind", kind)
	q.Add("date", date)
	q.Add("before_date", beforeDate)
	q.Add("after_date", afterDate)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("x-gladia-key", c.ApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var transcriptionListResponse TranscriptionListResponse
	err = json.Unmarshal(body, &transcriptionListResponse)
	if err != nil {
		return nil, err
	}

	return &transcriptionListResponse, nil
}
