package cmd

import (
	"fmt"
	"go/build"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

/*
For releases:

The version number and build date are set at build time via ldflags.
*/
var version = "dev"
var buildDate = time.Now().Format("2006-01-02 15:04:05")
var gitHash, _ = exec.Command("git", "rev-parse", "--short", "HEAD").Output()
var arch = build.Default.GOARCH
var kernel = build.Default.GOOS

// Flags
var versionShort = false

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information for yURL",
	Run: func(cmd *cobra.Command, args []string) {

		// This logic should be moved to Mage soon. Git only matters when building
		if strings.HasPrefix(version, "dev") || strings.HasPrefix(version, "SNAPSHOT") {

			if len(string(gitHash)) != 0 {
				hashString := string(gitHash)[:len(string(gitHash))-1]
				version = "dev-" + hashString
			}
		} else {
			buildDate = ""
			version = "v" + version
		}

		if versionShort {
			fmt.Print(version)
			return
		}

		fmt.Println("yURL")
		fmt.Println("Version: " + version)
		if buildDate != "" {
			fmt.Println("Date: " + buildDate)
		}
		fmt.Println("Platform: " + kernel + "/" + arch)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&versionShort, "short", "s", false, "Display just the actual version")
}
