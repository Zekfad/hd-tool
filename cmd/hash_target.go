package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Zekfad/hd-tool/game_data/reader"
	"github.com/Zekfad/hd-tool/hash_db"
	"github.com/spf13/cobra"
)

var targetCmd = &cobra.Command{
	Use:   "target [folder] [target]",
	Short: "Generate Hash DB target",
	Long:  `Scan folder for packages and contained hashes to form Hash DB target.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dirname := args[0]
		targetName := args[1]

		includePackageName, err := cmd.Flags().GetBool("package")
		if err != nil {
			fmt.Println("Failed to parse flag package")
			return
		}
		includeTypeName, err := cmd.Flags().GetBool("type")
		if err != nil {
			fmt.Println("Failed to parse flag type")
			return
		}
		includeFileName, err := cmd.Flags().GetBool("file")
		if err != nil {
			fmt.Println("Failed to parse flag file")
			return
		}

		if !(includePackageName || includeTypeName || includeFileName) {
			fmt.Println("No work to be done, flags are missing.")
			return
		}

		target := hash_db.HashDBTarget{}

		for fullPath, archive := range reader.ArchivesFromDirectory(dirname) {
			if includePackageName {
				filename := filepath.Base(fullPath)
				packageHash, err := strconv.ParseUint(filename, 16, 64)
				if err != nil {
					fmt.Printf("Warn: Package %s has non-hash name.", filename)
				} else {
					target[packageHash] = true
				}
			}

			for _, file := range archive.GetFiles() {
				if includeTypeName {
					target[uint64(file.GetType())] = true
				}
				if includeFileName {
					target[file.GetName()] = true
				}
			}
		}

		err = target.SaveToFile(targetName, true)
		if err != nil {
			fmt.Printf("Failed to save Hash DB Target: %s", err)
			return
		}
	},
}

var targetSortCmd = &cobra.Command{
	Use:   "sort [target]",
	Short: "Sort Hash DB target",
	Long:  `Sort Hash DB target and eliminate duplicates.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName := args[0]

		target, err := hash_db.TargetFromFile(targetName)
		if err != nil {
			fmt.Printf("Failed to load Hash DB Target: %s", err)
			return
		}
		err = target.SaveToFile(targetName, true)
		if err != nil {
			fmt.Printf("Failed to save Hash DB Target: %s", err)
			return
		}
	},
}

func init() {
	hashCmd.AddCommand(targetCmd)
	targetCmd.AddCommand(targetSortCmd)
	targetCmd.Flags().Bool("package", false, "include package names")
	targetCmd.Flags().Bool("type", false, "include type names")
	targetCmd.Flags().Bool("file", false, "include file names")
}
