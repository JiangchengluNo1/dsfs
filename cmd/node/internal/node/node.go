package node

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/mahaonan001/dsfs/cmd/node/internal/logic"
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
	data := make([]byte, 1024*1024*64)
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
	var file *os.File
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
			log.Println("accept file path: ", fc.Path)
			file, err = logic.OpenOrCreateFile(fc.Path)
			if err != nil {
				return err
			}
			defer file.Close()
		case *filetransfer.UploadFileRequest_Data:
			data := streamData.GetData()
			err := logic.WriteData(file, data)
			if err != nil {
				break
			}
		}
	}
	return stream.SendAndClose(&filetransfer.UploadFileResponse{Success: false})
}

func (n *Node) DeleteFile(ctx context.Context, in *filetransfer.DeleteFileRequest) (*filetransfer.DeleteFileResponse, error) {
	return nil, nil
}
