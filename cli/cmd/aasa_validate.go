package cmd

import (
	"fmt"

	"github.com/chayev/yurl/yurllib"

	"github.com/spf13/cobra"
)

// validateAASACmd represents the validate command for Apple App Site Association
var validateAASACmd = &cobra.Command{
	Use:   "validate <URL>",
	Short: "Validate your link against Apple's requirements",
	Run: func(cmd *cobra.Command, args []string) {
		output := yurllib.CheckAASADomain(args[0], "", "", true)

		for _, item := range output {
			fmt.Print(item)
		}

	},
}

func init() {
	aasaCmd.AddCommand(validateAASACmd)
}
