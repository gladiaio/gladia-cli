package types

// UploadResponse structure to capture response from file upload API
type UploadResponse struct {
	AudioURL string `json:"audio_url"` // Assuming the API returns a URL for the uploaded audio
}

// TranscriptionRequest structure for creating a transcription task
// TranscriptionRequest structure for creating a transcription task
type TranscriptionRequest struct {
	AudioURL            string               `json:"audio_url"`
	CallbackURL         string               `json:"callback_url,omitempty"`
	ContextPrompt       string               `json:"context_prompt,omitempty"`
	CustomMetadata      map[string]string    `json:"custom_metadata,omitempty"`
	CustomPrompts       bool                 `json:"custom_prompts,omitempty"`
	CustomPromptsConfig *CustomPromptsConfig `json:"custom_prompts_config,omitempty"`
	CustomVocabulary    []string             `json:"custom_vocabulary,omitempty"`
	DetectLanguage      bool                 `json:"detect_language,omitempty"`
	Diarization         bool                 `json:"diarization,omitempty"`
	DiarizationConfig   *DiarizationConfig   `json:"diarization_config,omitempty"`
	EnableCodeSwitching bool                 `json:"enable_code_switching,omitempty"`
	Language            Language             `json:"language,omitempty"`
	Subtitles           bool                 `json:"subtitles,omitempty"`
	SubtitlesConfig     *SubtitlesConfig     `json:"subtitles_config,omitempty"`
	Summarization       bool                 `json:"summarization,omitempty"`
	SummarizationConfig *SummarizationConfig `json:"summarization_config,omitempty"`
	Translation         bool                 `json:"translation,omitempty"`
	TranslationConfig   *TranslationConfig   `json:"translation_config,omitempty"`
}

type CustomPromptsConfig struct {
	Prompts []string `json:"prompts"`
}

type DiarizationConfig struct {
	MaxSpeakers      int `json:"max_speakers,omitempty"`
	MinSpeakers      int `json:"min_speakers,omitempty"`
	NumberOfSpeakers int `json:"number_of_speakers,omitempty"`
}

type SubtitlesConfig struct {
	Formats []SubtitleFormat `json:"formats"`
}

type SummarizationConfig struct {
	Type string `json:"type"`
}

type TranslationConfig struct {
	Model           string           `json:"model"`
	TargetLanguages []TargetLanguage `json:"target_languages"`
}

type SubtitleFormat string

// Constants for acceptable subtitle formats.
const (
	SubtitleFormatSRT SubtitleFormat = "srt"
	SubtitleFormatVTT SubtitleFormat = "vtt"
)
