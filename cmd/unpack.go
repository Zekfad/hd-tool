package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Zekfad/hd-tool/game_data"
	"github.com/Zekfad/hd-tool/game_data/reader"
	"github.com/Zekfad/hd-tool/hash_db"
	"github.com/spf13/cobra"
)

var unpackCmd = &cobra.Command{
	Use:   "unpack [archive] [target_dir]",
	Short: "Unpack game archive",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		archiveName := args[0]
		targetDirectory := args[1]

		dbName, err := cmd.Flags().GetString("hash-db")
		if err != nil {
			fmt.Println("Failed to parse hash db flag")
			return
		}

		unknown, err := cmd.Flags().GetBool("unknown")
		if err != nil {
			fmt.Println("Failed to parse unknown flag")
			return
		}

		var db = hash_db.HashDB{}
		if dbName != "" {
			db, err = hash_db.FromFile(dbName, false)
			if err != nil {
				fmt.Printf("Failed load hash db %s\n", err)
				return
			}
		}

		err = unpackFile(archiveName, targetDirectory, unknown, db)
		if err != nil {
			fmt.Print(err)
			return
		}
	},
}

var unpackAllCmd = &cobra.Command{
	Use:   "unpack-all [archives_dir] [target_dir]",
	Short: "Unpack game archives",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		archivesDirectory := args[0]
		targetDirectory := args[1]

		dbName, err := cmd.Flags().GetString("hash-db")
		if err != nil {
			fmt.Println("Failed to parse hash db flag")
			return
		}

		unknown, err := cmd.Flags().GetBool("unknown")
		if err != nil {
			fmt.Println("Failed to parse unknown flag")
			return
		}

		var db = hash_db.HashDB{}
		if dbName != "" {
			db, err = hash_db.FromFile(dbName, false)
			if err != nil {
				fmt.Printf("Failed load hash db %s\n", err)
				return
			}
		}

		entries, err := os.ReadDir(archivesDirectory)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range entries {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) == "" {
				fullPath := filepath.Join(archivesDirectory, file.Name())

				fmt.Printf("Trying to unpack candidate %s\n", fullPath)
				err = unpackFile(fullPath, targetDirectory, unknown, db)

				if err != nil {
					fmt.Printf("Failed: %s\n", err)
					continue
				}
				fmt.Printf("Unpacked %s!\n", fullPath)
			}
		}
	},
}

func unpackFile(name string, targetDirectory string, unknown bool, db hash_db.HashDB) error {
	if err := ensureDir(targetDirectory); err != nil {
		return errors.Join(
			fmt.Errorf("failed to create target directory"),
			err,
		)
	}
	archive, err := reader.ArchiveFromFile(name)
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to read archive"),
			err,
		)
	}

	fmt.Printf("Loaded archive of version: %#X\n", archive.GetVersion())
	for _, entry := range archive.GetFiles() {
		filename, dbHasName := db[uint64(entry.GetName())]
		if !dbHasName {
			filename = fmt.Sprintf("%016X", entry.GetName())
		}
		filePath := filepath.Join(targetDirectory, filename)
		switch entry.GetType() {
		case game_data.Type_lua:
			fmt.Printf("Found script %s ... ", filename)

			lua, err := game_data.LuaResourceFromBytes(entry.GetInlineBuffer())
			if err != nil {
				fmt.Printf("invalid: %s\n", err)
				break
			}
			fmt.Printf("valid (format: %#X)\n", lua.Format)

			if err := ensureDir(filepath.Dir(filePath)); err != nil {
				fmt.Printf("Warn: Failed to create target directory error: %s\n", err)
				break
			}

			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Printf("Warn: Failed to open target file: %s\n", err)
				break
			}
			defer file.Close()

			file.Write(lua.Data)
			break
		default:
			if !unknown {
				break
			}
			fmt.Printf("Found unknown file %s\n", filename)

			if err := ensureDir(filepath.Dir(filePath)); err != nil {
				fmt.Printf("Warn: Failed to create target directory error: %s\n", err)
				break
			}

			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				fmt.Printf("Warn: Failed to open target file: %s\n", err)
				break
			}
			defer file.Close()
			file.Write(entry.GetInlineBuffer())
		}
	}
	return nil
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func init() {
	rootCmd.AddCommand(unpackCmd)
	unpackCmd.Flags().String("hash-db", "", "Hash DB file.")
	unpackCmd.Flags().Bool("unknown", false, "Save raw buffer for unknown file formats")
	rootCmd.AddCommand(unpackAllCmd)
	unpackAllCmd.Flags().String("hash-db", "", "Hash DB file.")
	unpackAllCmd.Flags().Bool("unknown", false, "Save raw buffer for unknown file formats")
}
