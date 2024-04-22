package cmd

import (
	"fmt"
	"strconv"

	"github.com/Zekfad/hd-tool/game_data"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [folder] [target]",
	Short: "Search for package with a given type",
	Long:  `Scan folder for packages and files of given type.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dirname := args[0]
		target, err := strconv.ParseUint(args[1], 16, 64)
		if err != nil {
			fmt.Printf("failed to parse target: %s", err)
			return
		}

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println("Failed to parse all flag")
			return
		}

		targetType := game_data.TypeHash(target)

	archives_loop:
		for path, archive := range game_data.ArchivesFromDirectory(dirname) {
			fmt.Println("Checking:", path)
			for i, file := range archive.Files {
				if targetType == file.Type {
					fmt.Printf("Found match in archive: %s file %d\n", path, i)
					if !all {
						break archives_loop
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Bool("all", false, "find all matches")
}
