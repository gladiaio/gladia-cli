package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// version is set at build time via ldflags: -X main.version=<tag>
var version = "dev"

var rootGladiaKey string

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "gladia",
		Short:   "Gladia speech-to-text CLI",
		Long:    "Transcribe audio files and URLs with the Gladia API.",
		Version: version,
	}

	rootCmd.PersistentFlags().StringVar(&rootGladiaKey, "gladia-key", "", "Gladia API key (used when GLADIA_API_KEY and ~/.gladia are unset)")

	rootCmd.AddCommand(newTranscribeCmd())
	rootCmd.AddCommand(newAuthCmd())
	rootCmd.AddCommand(newLanguagesCmd())

	return rootCmd
}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
