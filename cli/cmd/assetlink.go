package cmd

import (
	"github.com/spf13/cobra"
)

var (
	assetlinkCmd = &cobra.Command{
		Use:   "assetlink",
		Short: "Command for Android Asset Link utils.",
	}
)

func init() {
	rootCmd.AddCommand(assetlinkCmd)
}
