package cmd

import (
	"github.com/spf13/cobra"
)

var (
	aasaCmd = &cobra.Command{
		Use:   "aasa",
		Short: "Command for Apple App Site Association utils.",
	}
)

func init() {
	rootCmd.AddCommand(aasaCmd)
}
