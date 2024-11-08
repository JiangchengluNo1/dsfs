package logic

import (
	"errors"
	"fmt"
	"os"

	filetransfer "github.com/mahaonan001/dsfs/proto"
)

// GetFile reads a file from the filesystem, returning the contents as a byte slice.
func GetFile(path string, stream *filetransfer.FileTransfer_GetFileServer) error {
	shas := FileHolder.GetFile(path)
	fmt.Println(shas)
	if len(shas) == 0 {
		return errors.New("file not found")
	}
	buf := make([]byte, 1024*1024*64)
	for _, sha := range shas {
		path := fmt.Sprintf("%s/%x", fileDir, sha)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if len(data) != len(buf) {
			return errors.New("data length not equal")
		}
		if err := (*stream).Send(&filetransfer.FileChunk{Data: data}); err != nil {
			return err
		}
	}
	return nil
}
