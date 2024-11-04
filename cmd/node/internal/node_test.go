package main

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	filetransfer "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRPC(t *testing.T) {
	t.Log("test rpc")
	conn, err := grpc.NewClient("localhost:5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := filetransfer.NewFileTransferClient(conn)

	file, err := os.Open("./file/test.meta")
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	defer file.Close()

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		log.Fatalf("could not upload file: %v", err)
	}
	if err := stream.Send(&filetransfer.UploadFileRequest{
		Payload: &filetransfer.UploadFileRequest_Fm{
			Fm: &filetransfer.FileMeta{Path: "new.meta", Index: 0},
		},
	}); err != nil {
		log.Fatalf("could not send file meta: %v", err)
	}
	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not read file: %v", err)
		}

		if err := stream.Send(&filetransfer.UploadFileRequest{
			Payload: &filetransfer.UploadFileRequest_Data{Data: buf[:n]},
		}); err != nil {
			log.Fatalf("could not send chunk: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not receive response: %v", err)
	}

	log.Printf("Upload response: %v", res.Success)
}
