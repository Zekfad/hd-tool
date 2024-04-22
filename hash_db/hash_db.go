package hash_db

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/maruel/natural"
)

type HashDB map[uint64]string

func FromFile(name string, checkHashes bool) (HashDB, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()

	hashDb := HashDB{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// skip empty lines and comments starting with #
		if line == "" || line[0] == '#' {
			continue
		}
		key, value, valid := strings.Cut(line, " ")
		if !valid {
			return nil, fmt.Errorf("invalid hash db line: %s", line)
		}
		hash, err := strconv.ParseUint(key, 16, 64)
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to parse key: %s", key),
				err,
			)
		}
		if checkHashes {
			if expected := Hash(value); expected != hash {
				return nil, fmt.Errorf(
					"invalid hash expected %016x got %016x value: %s",
					expected,
					hash,
					value,
				)
			}
		}
		hashDb[hash] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to parse hash db"),
			err,
		)
	}

	return hashDb, nil
}

func (db HashDB) SaveToFile(name string, sortKeys bool) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()
	if sortKeys {
		var keys []uint64
		for key := range db {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i, j int) bool {
			return natural.Less(db[keys[i]], db[keys[j]])
		})

		for _, key := range keys {
			value := db[key]
			_, err := file.WriteString(
				fmt.Sprintf("%016X %s\n", key, value),
			)
			if err != nil {
				return errors.Join(
					fmt.Errorf("failed write hash db"),
					err,
				)
			}
		}
	} else {
		for key, value := range db {
			_, err := file.WriteString(
				fmt.Sprintf("%016X %s\n", key, value),
			)
			if err != nil {
				return errors.Join(
					fmt.Errorf("failed write hash db"),
					err,
				)
			}
		}
	}
	return nil
}

func (db HashDB) AddHash(value string) {
	db[Hash(value)] = value
}

func (db HashDB) AddHashesFromFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		// cut off \n and \r\n
		if len := len(line); len > 1 {
			if line[len-2] == '\r' {
				line = line[:len-2]
			} else {
				line = line[:len-1]
			}
		}

		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to scan file"),
				err,
			)
		}
		// skip empty lines
		if line == "" || line == "\n" || line == "\r\n" || line == "\r" {
			continue
		}
		db.AddHash(line)
	}
	return nil
}
