package cmd

import (
	"fmt"

	"github.com/chayev/yurl/yurllib"

	"github.com/spf13/cobra"
)

// validateAssetLinkCmd represents the validate command for ASset Links
var validateAssetLinkCmd = &cobra.Command{
	Use:   "validate <URL>",
	Short: "Validate your link against Android's requirements",
	Run: func(cmd *cobra.Command, args []string) {
		output := yurllib.CheckAssetLinkDomain(args[0], "", "")

		for _, item := range output {
			fmt.Print(item)
		}

	},
}

func init() {
	assetlinkCmd.AddCommand(validateAssetLinkCmd)
}
