package format

import (
	"encoding/json"
	"fmt"
	transcribe "gladia/api/transcribe"
	"os"

	"github.com/olekukonko/tablewriter"
)

func PrintTXTTranscription(response transcribe.TranscriptionResult) {
	println()
	fmt.Println(response.Result.Transcription.FullTranscript)
}

func PrintTXTDiarizedTranscription(response transcribe.TranscriptionResult) {
	for _, utterance := range response.Result.Transcription.Utterances {
		// Print the speaker label and the transcription text
		fmt.Printf("Speaker %d: %s\n", utterance.Speaker, utterance.Text)
	}
}

func PrintTableTranscription(response transcribe.TranscriptionResult) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"TIME BEGIN", "TIME END", "LANGUAGE", "SPEAKER", "TRANSCRIPTION"})
	table.SetBorder(true) // Enable borders around the table
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true) // Enable line after the header

	utterances := response.Result.Transcription.Utterances

	for _, transcription := range utterances {
		table.Append([]string{
			fmt.Sprintf("%.6f", transcription.Start),
			fmt.Sprintf("%.6f", transcription.End),
			transcription.Language,
			fmt.Sprintf("%d", transcription.Speaker),
			transcription.Text,
		})
	}

	table.Render() // Render the table to os.Stdout
}

func PrintJSONTranscription(response transcribe.TranscriptionResult) {
	jsonData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))

}

func PrintJSONSimplifiedTranscription(response transcribe.TranscriptionResult) {
	// Check if there are any utterances to avoid index out of range errors
	if len(response.Result.Transcription.Utterances) == 0 {
		fmt.Println("No transcriptions available.")
		return
	}

	// Define a slice to hold all simplified transcriptions
	var simplifiedTranscriptions []struct {
		TimeBegin     string `json:"time_begin"`
		TimeEnd       string `json:"time_end"`
		Language      string `json:"language"`
		Speaker       string `json:"speaker"`
		Transcription string `json:"transcription"`
	}

	// Iterate over all utterances and add them to the slice
	for _, utterance := range response.Result.Transcription.Utterances {
		simplifiedTranscriptions = append(simplifiedTranscriptions, struct {
			TimeBegin     string `json:"time_begin"`
			TimeEnd       string `json:"time_end"`
			Language      string `json:"language"`
			Speaker       string `json:"speaker"`
			Transcription string `json:"transcription"`
		}{
			TimeBegin:     fmt.Sprintf("%.6f", utterance.Start),
			TimeEnd:       fmt.Sprintf("%.6f", utterance.End),
			Language:      utterance.Language,
			Speaker:       fmt.Sprintf("%d", utterance.Speaker),
			Transcription: utterance.Text,
		})
	}

	// Marshal the slice of simplified transcriptions into JSON
	jsonData, err := json.MarshalIndent(simplifiedTranscriptions, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func PrintCSVTranscription(response transcribe.TranscriptionResult) {
	// Print the CSV header
	fmt.Println("time_begin, time_end, language, speaker, transcription")

	// Iterate over each utterance and print its details in CSV format
	for _, utterance := range response.Result.Transcription.Utterances {
		fmt.Printf("%.6f, %.6f, %s, %d, \"%s\"\n",
			utterance.Start,
			utterance.End,
			utterance.Language,
			utterance.Speaker,
			utterance.Text)
	}
}

func PrintSRTTranscription(response transcribe.TranscriptionResult) {
	for i, utterance := range response.Result.Transcription.Utterances {
		// Convert start and end times from seconds to SRT time format
		startTime := secondsToSRTTimeFormat(utterance.Start)
		endTime := secondsToSRTTimeFormat(utterance.End)

		// Print the SRT block
		fmt.Printf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, utterance.Text)
	}
}

// Helper function to convert seconds to SRT time format (HH:MM:SS,MS)
func secondsToSRTTimeFormat(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := int(seconds) / 60 % 60
	secs := int(seconds) % 60
	milliseconds := int(seconds*1000) % 1000

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, milliseconds)
}

func PrintSRTDiarizedTranscription(response transcribe.TranscriptionResult) {
	for i, utterance := range response.Result.Transcription.Utterances {
		// Convert start and end times from seconds to SRT time format
		startTime := secondsToSRTTimeFormat(utterance.Start)
		endTime := secondsToSRTTimeFormat(utterance.End)

		// Print the SRT block with speaker diarization
		fmt.Printf("%d\n%s --> %s\nSpeaker %d: %s\n\n", i+1, startTime, endTime, utterance.Speaker, utterance.Text)
	}
}

func PrintVTTTranscription(response transcribe.TranscriptionResult) {
	// Print the WebVTT file header
	fmt.Println("WEBVTT")
	fmt.Println()

	for i, utterance := range response.Result.Transcription.Utterances {
		// Convert start and end times from seconds to WebVTT time format
		startTime := secondsToVTTTimeFormat(utterance.Start)
		endTime := secondsToVTTTimeFormat(utterance.End)

		// Print the cue identifier, time range, and transcription text
		fmt.Printf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, utterance.Text)
	}
}

// Helper function to convert seconds to WebVTT time format (HH:MM:SS.MMM)
func secondsToVTTTimeFormat(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := int(seconds) / 60 % 60
	secs := int(seconds) % 60
	milliseconds := int(seconds*1000) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, milliseconds)
}

func PrintVTTDiarizedTranscription(response transcribe.TranscriptionResult) {
	// Print the VTT header
	fmt.Println("WEBVTT")
	fmt.Println()

	for i, utterance := range response.Result.Transcription.Utterances {
		// Convert start and end times from seconds to WebVTT time format
		startTime := secondsToVTTTimeFormat(utterance.Start)
		endTime := secondsToVTTTimeFormat(utterance.End)

		// Print the cue identifier, time range, speaker label, and transcription text
		fmt.Printf("%d\n%s --> %s\nSpeaker %d: %s\n\n", i+1, startTime, endTime, utterance.Speaker, utterance.Text)
	}
}

func PrintSummarization(response transcribe.TranscriptionResult) {
	if response.Result.Summarization.Results != nil {
		println()
		fmt.Println(*response.Result.Summarization.Results)
	} else {
		fmt.Println("No summarization results available.")
	}
}
