package logic

import (
	"fmt"
	"os"
)

type UpdateLog struct {
	From     int
	To       int
	FileMeta [32]byte
}

var (
	Updl = make(map[[32]byte]UpdateLog)
)

func UpdateFile(filemeta [32]byte, up_data []byte, from, to int) error {
	// get the file path
	sha := GenerateSHA256(up_data)
	f, err := os.OpenFile(fmt.Sprintf("%s/%x", fileDir, sha), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(up_data)
	if err != nil {
		return err
	}
	FileHolder.AppendFile(fmt.Sprintf("%x", filemeta), sha)
	if n != len(up_data) {
		return os.ErrDeadlineExceeded
	}
	Updl[filemeta] = UpdateLog{
		From:     from,
		To:       to,
		FileMeta: sha,
	}
	return nil
}
