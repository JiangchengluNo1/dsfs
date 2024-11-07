package logic

import (
	"crypto/sha256"
	"fmt"
	"os"
)

const fileDir = "./file"

func OpenOrCreateFile(path string) (*os.File, error) {
	fmt.Println("receive file path:", path)
	file, err := os.OpenFile(fileDir+"/"+path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func WriteData(file *os.File, data []byte) error {
	n, err := file.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return os.ErrDeadlineExceeded
	}
	return nil
}

func GenerateSHA256(data []byte) [32]byte {
	sum := sha256.Sum256(data)
	return sum
}

func CheckSumExisted(sum [32]byte) bool {
	// check if the file with the same checksum exists
	return false
}
