package main

import (
	"fmt"
	"net"

	"github.com/mahaonan001/dsfs/cmd/node/internal/logic"
	"github.com/mahaonan001/dsfs/cmd/node/internal/node"
	filetransfer "github.com/mahaonan001/dsfs/proto/transfer"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*128), grpc.MaxRecvMsgSize(1024*1024*128))
	filetransfer.RegisterFileTransferServer(s, &node.Node{})
	fmt.Println("Node start at 5001")
	go logic.FileHolder.KeepFlushBuffer()
	defer logic.FileHolder.Close()
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func init() {
	logic.Init()
	files := make([]string, 0, logic.FileHolder.Getlen())
	for k := range logic.FileHolder.GetMap() {
		files = append(files, k)
	}
	node.WakeUp(files)
	go node.Healthing(1)
}
