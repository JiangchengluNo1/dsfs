package node

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/mahaonan001/dsfs/cmd/node/internal/logic"
	filetransfer "github.com/mahaonan001/dsfs/proto"
)

type (
	NodeModel interface {
		filetransfer.UnimplementedFileTransferServer
	}
	CustomNodeModel interface {
		filetransfer.FileTransferServer
	}
)

type Node struct {
	filetransfer.UnimplementedFileTransferServer
}

var Path string

// *filetransfer.GetFileRequest, grpc.ServerStreamingServer[filetransfer.FileChunk]
func (n *Node) GetFile(req *filetransfer.GetFileRequest, stream filetransfer.FileTransfer_GetFileServer) error {

	err := logic.GetFile(req.Path, &stream)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) UploadFile(stream filetransfer.FileTransfer_UploadFileServer) error {
	for {
		streamData, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Upload file success\n")
			return stream.SendAndClose(&filetransfer.UploadFileResponse{Success: true})
		}
		if err != nil {
			break
		}
		switch streamData.Payload.(type) {
		case *filetransfer.UploadFileRequest_Fm:
			fc := streamData.GetFm()
			Path = fc.Path
		case *filetransfer.UploadFileRequest_Data:
			data := streamData.GetData()
			sum := logic.GenerateSHA256(data)
			exist := logic.CheckSumExisted(Path, sum)
			if exist {
				/*软链接path与sha256对应的文件*/
				continue
			}
			sha := fmt.Sprintf("%x", sum)
			fmt.Printf("file %s not exists\n", sha)
			logic.FileHolder.AppendFile(Path, sum)
			_, err := logic.WriteData(sum, data) //add data to file
			if err != nil {
				log.Println("write data to file error:", err)
				break
			}
		}
	}
	return stream.SendAndClose(&filetransfer.UploadFileResponse{Success: false})
}

func (n *Node) DeleteFile(ctx context.Context, in *filetransfer.DeleteFileRequest) (*filetransfer.DeleteFileResponse, error) {
	return nil, nil
}
