package logic

import (
	"crypto/md5"
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
	/*TODO:check the data md5 value*/
	n, err := file.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return os.ErrDeadlineExceeded
	}
	return nil
}

func generateMd5(data []byte) [16]byte {
	md := md5.Sum(data)
	return md
}
