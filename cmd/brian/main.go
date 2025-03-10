package brian

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "brian",
	Short: "A simple CLI tool that responds with hello from Brian",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from Brian")
	},
}

func Execute() error {
	return rootCmd.Execute()
}