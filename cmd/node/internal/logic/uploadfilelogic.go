package logic

import (
	"crypto/sha256"
	"fmt"
	"os"
)

const (
	fileDir = "../cmd/node/internal/file"
)

func WriteData(sha [32]byte, data []byte) ([32]byte, error) {
	shaPath := fmt.Sprintf("%x", sha)
	file, err := os.OpenFile(fileDir+"/"+shaPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return [32]byte{}, err
	}
	defer file.Close()
	n, err := file.Write(data)
	if err != nil {
		return [32]byte{}, err
	}
	if n != len(data) {
		return [32]byte{}, os.ErrDeadlineExceeded
	}
	return [32]byte{}, nil
}

func GenerateSHA256(data []byte) [32]byte {
	sum := sha256.Sum256(data)
	return sum
}

func CheckSumExisted(path string, sum [32]byte) bool {
	// check if the file with the same checksum exists
	shas, ok := FileHolder.m[path]
	if ok {
		for _, sha := range shas {
			if sha == sum {
				return true
			}
		}
	}
	return false
}
