package cmd

import (
	"fmt"

	"github.com/Zekfad/hd-tool/hash_db"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db [db_file] [source_file...]",
	Short: "Update Hash DB",
	Long: `Hash DB utility

This utility will help you to update Hash DB.
Source files are read by lines, without trimming of trailing spaces.
Target file contains include list of hashes.

All hashes are written and read as base-16 (hex) 64 bit unsigned integers.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName, err := cmd.Flags().GetString("target")
		if err != nil {
			fmt.Println("Failed to parse flag target")
			return
		}
		check, err := cmd.Flags().GetBool("check")
		if err != nil {
			fmt.Println("Failed to parse flag check")
			return
		}
		sort, err := cmd.Flags().GetBool("sort")
		if err != nil {
			fmt.Println("Failed to parse flag sort")
			return
		}
		dbName := args[0]
		sources := args[1:]
		hasSources := len(sources) > 0
		hasTarget := targetName != ""

		db, err := hash_db.FromFile(dbName, check)
		if err != nil {
			fmt.Printf("Failed to load Hash DB: %s\n", err)
			return
		}
		fmt.Printf("Hash DB loaded successfully (%d entries)\n", len(db))

		if hasSources {
			for _, file := range sources {
				fmt.Printf("Processing %s ... ", file)

				err = db.AddHashesFromFile(file)
				if err != nil {
					fmt.Printf("failed: %s\n", err)
				} else {
					fmt.Printf("done\n")
				}
			}
		}

		if hasTarget {
			target, err := hash_db.TargetFromFile(targetName)
			if err != nil {
				fmt.Printf("Failed to load target: %s\n", err)
				return
			}
			fmt.Printf("Successfully load target %s (%d entries)\n", targetName, len(target))

			for hash := range db {
				if !target[hash] {
					delete(db, hash)
				}
			}
		}

		err = db.SaveToFile(dbName, sort)
		if err != nil {
			fmt.Printf("Failed to save Hash DB: %s\n", err)
			return
		}
		fmt.Printf("Hash DB saved successfully (%d entries)\n", len(db))
	},
}

func init() {
	hashCmd.AddCommand(dbCmd)
	dbCmd.Flags().Bool("check", false, "Check Hash DB values on load")
	dbCmd.Flags().Bool("sort", true, "Sort Hash DB on save")
	dbCmd.Flags().String("target", "", "Target file - apply include filter to Hash DB.")
}
