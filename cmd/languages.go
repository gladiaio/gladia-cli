package main

import (
	"github.com/gladiaio/gladia-cli/pkg/client/types"
	"github.com/spf13/cobra"
)

func newLanguagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "languages",
		Short: "List supported transcription language codes",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := types.DisplayAllInputLanguagesNames()
			return err
		},
	}
}
