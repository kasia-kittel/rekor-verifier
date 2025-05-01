package utils

import (
	"crypto/sha256"
	"hash"
	"io"
	"os"

	"github.com/pkg/errors"
)

func CheckPathToFile(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, "Path doesn't exist")
		} else {
			return errors.Wrap(err, "Error accessing path")
		}
	}
	if !fi.Mode().IsRegular() {
		return errors.New("Path to directory")
	}

	return nil
}


func CalculateSHA256(file *os.File) (hash.Hash, error) {
	sha := sha256.New()

	if _, err := io.Copy(sha, file); err != nil {
		return nil, errors.Wrap(err, "Error calculating SHA265")
	}

	return sha, nil
}