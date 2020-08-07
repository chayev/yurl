package cmd

import (
	"fmt"

	"github.com/chayev/yurl/yurllib"

	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <URL>",
	Short: "Validate your link against Apple's requirements",
	Run: func(cmd *cobra.Command, args []string) {
		output := yurllib.CheckDomain(args[0], "", "", true)

		for _, item := range output {
			fmt.Print(item)
		}

	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
