package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage Gladia API credentials",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "set [api-key]",
		Short: "Save your API key to ~/.gladia",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := SaveGladiaKeyToFile(args[0]); err != nil {
				return fmt.Errorf("save API key: %w", err)
			}
			return nil
		},
	})

	return cmd
}
