package main

import (
	"io"

	filetransfer "github.com/mahaonan001/dsfs/proto"
)

type FileServer struct {
	filetransfer.UnimplementedFileServiceServer
}

// 实现文件上传
func (f *FileServer) UploadFile(stream filetransfer.FileService_UploadServer) error {
	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {

		} else if err != nil {
			return err
		}
	}
}

func main() {

}
