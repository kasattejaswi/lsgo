package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lsgo",
	Short: "Go implementation of ls command in bash",
	Long: `This is the go implementation of some features of ls command
For example:
  lsgo list
  lsgo mod`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
