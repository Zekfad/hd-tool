package reader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"

	"github.com/Zekfad/hd-tool/game_data"
	"github.com/Zekfad/hd-tool/game_data/hd1"
	"github.com/Zekfad/hd-tool/game_data/hd2"
	"github.com/ghostiam/binstruct"
)

func ArchiveFromBytes(data []byte) (game_data.Archive, error) {
	reader := binstruct.NewReaderFromBytes(data, binary.LittleEndian, false)
	version, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	reader.Seek(0, io.SeekStart)
	switch game_data.ArchiveVersion(version) {
	case game_data.ArchiveVersionHD2:
		var archive hd2.Archive
		err = reader.Unmarshal(&archive)
		if err != nil {
			return nil, err
		}
		return archive, nil
	case game_data.ArchiveVersionHD1:
		var archive hd1.Archive
		err = reader.Unmarshal(&archive)
		if err != nil {
			return nil, err
		}
		return archive, nil
	default:
		return nil, fmt.Errorf("failed to decode archive of unknown version %#X", version)
	}
}

func ArchiveFromFile(name string) (game_data.Archive, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()

	buffer, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to read file"),
			err,
		)
	}

	archive, err := ArchiveFromBytes(buffer)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to parse file"),
			err,
		)
	}
	return archive, nil
}

func ArchivesFromDirectory(dirname string) iter.Seq2[string, game_data.Archive] {
	return func(yield func(path string, archive game_data.Archive) bool) {
		entries, err := os.ReadDir(dirname)
		if err != nil {
			fmt.Printf("Failed to read directory: %s", err)
			return
		}

		for _, file := range entries {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) == "" {
				fullPath := filepath.Join(dirname, file.Name())

				// fmt.Printf("Candidate %s ... ", fullPath)
				archive, err := ArchiveFromFile(fullPath)
				if err != nil {
					// fmt.Printf("not a package: %s\n", err)
					continue
				}
				// fmt.Printf("is a valid package!\n")
				if !yield(fullPath, archive) {
					return
				}
			}
		}
	}
}
