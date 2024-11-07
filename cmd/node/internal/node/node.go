package node

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mahaonan001/dsfs/cmd/node/internal/logic"
	"github.com/mahaonan001/dsfs/cmd/node/tys"
	filetransfer "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
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

// *filetransfer.GetFileRequest, grpc.ServerStreamingServer[filetransfer.FileChunk]
func (n *Node) GetFile(req *filetransfer.GetFileRequest, stream grpc.ServerStreamingServer[filetransfer.FileChunk]) error {
	var (
		err  error
		file *os.File
		ln   int
	)
	file, err = logic.GetFile(req.Path)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make([]byte, tys.MB*64)
	for {
		ln, err = file.Read(data)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = stream.Send(&filetransfer.FileChunk{Data: data[:ln]})
		if err != nil {
			break
		}
	}
	return err
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
		var path string
		switch streamData.Payload.(type) {
		case *filetransfer.UploadFileRequest_Fm:
			fc := streamData.GetFm()
			path = fc.Path
		case *filetransfer.UploadFileRequest_Data:
			data := streamData.GetData()
			sum := logic.GenerateSHA256(data)
			exist := logic.CheckSumExisted(sum)
			if exist {
				/*软链接path与sha256对应的文件*/
				fmt.Println("file already exists")
				continue
			}
			logic.FileHolder.AppendFile(path, sum)
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
