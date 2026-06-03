package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootGladiaKey string

func main() {
	rootCmd := &cobra.Command{
		Use:   "gladia",
		Short: "Gladia speech-to-text CLI",
		Long:  "Transcribe audio files and URLs with the Gladia API.",
	}

	rootCmd.PersistentFlags().StringVar(&rootGladiaKey, "gladia-key", "", "Gladia API key (used when GLADIA_API_KEY and ~/.gladia are unset)")

	rootCmd.AddCommand(newTranscribeCmd())
	rootCmd.AddCommand(newAuthCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
