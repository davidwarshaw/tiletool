package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the git hash version of tiletool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tiletool version %s\n", Version)
	},
}
