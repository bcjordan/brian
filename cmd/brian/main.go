package brian

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

// Version is the current version of the CLI
const Version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:   "brian",
	Short: "A simple CLI tool that responds with hello from Brian",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from Brian")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("brian version %s\n", Version)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and install updates",
	Run: func(cmd *cobra.Command, args []string) {
		doSelfUpdate()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(updateCmd)
}

func doSelfUpdate() {
	v := semver.MustParse(Version)
	latest, found, err := selfupdate.DetectLatest("bcjordan/brian")
	if err != nil {
		log.Println("Error occurred while checking for updates:", err)
		return
	}

	if !found {
		log.Println("No release found on GitHub. Are you running a local build?")
		return
	}

	if latest.Version.LTE(v) {
		fmt.Println("Current version", Version, "is the latest")
		return
	}

	fmt.Printf("New version v%s available\n", latest.Version)
	fmt.Printf("Release notes:\n%s\n", latest.ReleaseNotes)

	fmt.Print("Do you want to update? (y/n): ")
	input := ""
	fmt.Scanln(&input)
	if input != "y" && input != "Y" {
		return
	}

	exe, err := os.Executable()
	if err != nil {
		log.Println("Could not locate executable path:", err)
		return
	}

	if runtime.GOOS == "windows" {
		// Windows won't allow replacing a running executable
		// You'd need a more complex solution with a separate updater
		log.Println("Self-update on Windows requires closing the application first")
		log.Println("Please run 'go install github.com/bcjordan/brian@latest' to update")
		time.Sleep(2 * time.Second)
		return
	}

	fmt.Println("Updating...")
	err = selfupdate.UpdateTo(latest.AssetURL, exe)
	if err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	
	fmt.Println("Successfully updated to version", latest.Version)
}

func Execute() error {
	// Check for updates in the background if running the main command
	// Skip for version and update commands
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] != "version" && os.Args[1] != "update") {
		go func() {
			v := semver.MustParse(Version)
			latest, found, err := selfupdate.DetectLatest("bcjordan/brian")
			if err != nil || !found || latest.Version.LTE(v) {
				return
			}
			fmt.Fprintf(os.Stderr, "\n\nNew version v%s available! Run 'brian update' to upgrade.\n\n", latest.Version)
		}()
	}

	return rootCmd.Execute()
}