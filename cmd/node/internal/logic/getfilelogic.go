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
	if len(shas) == 0 {
		return errors.New("file not found")
	}
	for _, sha := range shas {
		path := fmt.Sprintf("%s/%x", fileDir, sha)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := (*stream).Send(&filetransfer.FileChunk{Data: data}); err != nil {
			return err
		}
	}
	return nil
}
