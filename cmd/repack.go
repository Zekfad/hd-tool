package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Zekfad/hd-tool/game_data"
	"github.com/Zekfad/hd-tool/game_data/hd1"
	"github.com/Zekfad/hd-tool/game_data/reader"
	"github.com/Zekfad/hd-tool/game_data/writer/writer_hd1"
	"github.com/Zekfad/hd-tool/hash_db"
	"github.com/spf13/cobra"
)

var repackCmd = &cobra.Command{
	Use:   "repack [original_archive] [new_archive] [patch_directory]",
	Short: "Patch game archive (only HD1)",
	Long:  `Patch HD1 archive to replace Lua scripts`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		original := args[0]
		new := args[1]
		patchDirectory := args[2]

		dbName, err := cmd.Flags().GetString("hash-db")
		if err != nil {
			fmt.Println("Failed to parse hash db flag")
			return
		}

		compiler, err := cmd.Flags().GetString("compiler")
		if err != nil {
			fmt.Println("Failed to parse compiler flag")
			return
		}
		if compiler == "" {
			fmt.Println("Compiler is required")
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

		err = repackFile(original, new, compiler, patchDirectory, db)
		if err != nil {
			fmt.Print(err)
			return
		}
	},
}

func repackFile(
	original string,
	new string,
	compiler string,
	patchDirectory string,
	db hash_db.HashDB,
) error {
	if err := ensureDir(patchDirectory); err != nil {
		return errors.Join(
			fmt.Errorf("failed to create patch directory"),
			err,
		)
	}
	src, err := reader.ArchiveFromFile(original)
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to read original archive"),
			err,
		)
	}
	archive, ok := src.(hd1.Archive)
	if !ok {
		return fmt.Errorf("unsupported archive")
	}

	for _, entry := range archive.Unpacked.Files {
		if entry.GetType() != game_data.Type_lua {
			continue
		}

		filename, dbHasName := db[uint64(entry.GetName())]
		if !dbHasName {
			filename = fmt.Sprintf("%016X", entry.GetName())
		}
		filePath := filepath.Join(patchDirectory, filename+".lua")

		// fmt.Printf("Searching for patch %s\n", filePath)
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			continue
		}

		fmt.Printf("Compiling patch script %s ... ", filePath)

		data, err := compileLuaJIT(compiler, filePath)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		fmt.Printf("done ... ")
		lua, _ := game_data.LuaResourceFromBytes(entry.GetInlineBuffer())
		lua.Data = data
		patchedResource, _ := lua.ToBytes()
		entry.VariantBuffers[0] = patchedResource

		fmt.Printf("patched %s!\n", filename)
	}

	file, err := os.OpenFile(new, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()
	return writer_hd1.WriteArchive(archive, file)
}

func compileLuaJIT(luajit string, script string) ([]byte, error) {
	dir := filepath.Dir(luajit)
	tempFile, err := os.CreateTemp("", "lua_bytecode_*")
	if err != nil {
		return nil, err
	}
	temp := tempFile.Name()
	defer os.Remove(temp)
	defer tempFile.Close()
	cmd := exec.Command(luajit, "-b", script, temp)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("LUA_PATH=%s", filepath.Join(dir, "?.lua")),
	)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return os.ReadFile(temp)
}

func init() {
	rootCmd.AddCommand(repackCmd)
	repackCmd.Flags().String("hash-db", "", "Hash DB file.")
	repackCmd.Flags().String("compiler", "", "Required: Path to LuaJIT 2.0.3")
}
