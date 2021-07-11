package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print out the version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lsgo           -   Version 1.0")
		fmt.Println("Developed by   -   Tejaswi Kasat")
		fmt.Println("Repository     -   https://github.com/kasattejaswi/lsgo.git")
		fmt.Println("Open for contribution")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
