package hash_db

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
)

type HashDBTarget map[uint64]bool

func TargetFromFile(name string) (HashDBTarget, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("failed to open file"),
			err,
		)
	}
	defer file.Close()

	hashDbTarget := HashDBTarget{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key := scanner.Text()
		// skip empty lines and comments starting with #
		if key == "" || key[0] == '#' {
			continue
		}
		hash, err := strconv.ParseUint(key, 16, 64)
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to parse key: %s", key),
				err,
			)
		}
		hashDbTarget[hash] = true
	}
	return hashDbTarget, nil
}

func (target HashDBTarget) SaveToFile(name string, sortKeys bool) error {
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
		for key := range target {
			keys = append(keys, key)
		}
		slices.SortFunc(keys, func(a uint64, b uint64) int {
			if a > b {
				return 1
			} else if a == b {
				return 0
			} else {
				return -1
			}
		})
		for _, key := range keys {
			_, err := file.WriteString(
				fmt.Sprintf("%016X\n", key),
			)
			if err != nil {
				return errors.Join(
					fmt.Errorf("failed write hash db"),
					err,
				)
			}
		}
	} else {
		for key := range target {
			_, err := file.WriteString(
				fmt.Sprintf("%016X\n", key),
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
