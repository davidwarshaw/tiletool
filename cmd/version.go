package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string
var Commit string
var Date string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version of tiletool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tiletool version %s, built %s\n", Version, Date)
	},
}
