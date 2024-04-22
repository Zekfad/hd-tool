package cmd

import (
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash DB and Hash DB Target utilities",
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
