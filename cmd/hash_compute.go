package cmd

import (
	"fmt"

	"github.com/Zekfad/hd-tool/hash_db"
	"github.com/spf13/cobra"
)

var computeCmd = &cobra.Command{
	Use:   "compute [string...]",
	Short: "Compute hash of value",
	Long:  "Computes entires suitable for Hash DB",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, value := range args {
			hash := hash_db.Hash(value)
			fmt.Printf("%016X %s\n", hash, value)
		}
	},
}

func init() {
	hashCmd.AddCommand(computeCmd)
}
