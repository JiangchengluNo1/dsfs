package main

import (
	"fmt"
	"net"

	"github.com/mahaonan001/dsfs/cmd/node/internal/logic"
	"github.com/mahaonan001/dsfs/cmd/node/internal/node"
	filetransfer "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
)

// 姐姐网恋吗？ ❤️‍🔥o.O💞
func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*128), grpc.MaxRecvMsgSize(1024*1024*128))
	filetransfer.RegisterFileTransferServer(s, &node.Node{})
	fmt.Println("Node start at 5001")
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
	defer logic.FileHolder.Close()
}

// 姐姐一个人写代码孤单么？ 💓O.o 💖
// 姐姐加一下我的绿泡泡吧！ 💚🧊💚
