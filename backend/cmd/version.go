package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Lazy Rabbit Secretary",
	Long:  `All software has versions. This is Lazy Rabbit Secretary`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Lazy Rabbit Secretary v0.1 -- Walter Fan")
	},
}
