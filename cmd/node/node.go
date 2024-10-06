package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/mahaonan001/dsfs/cmd/node/model"
	pb "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var MasterOnline = false

type NodeStorageServer struct {
	model.NodeServer
	pb.UnimplementedNodeServerServer
}

func (n *NodeStorageServer) TransportFileServer(ctx context.Context, sf *pb.SpiltedFile) (*pb.SplitedFileRes, error) {
	// for{

	// }
	// n.AddFile()
	return &pb.SplitedFileRes{Status: nil}, nil
}

func main() {
	go client2Master()
}

func client2Master() {
	// 设置与 Master 的 gRPC 连接
	conn, err := grpc.NewClient("localhost:45678", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}

	// 创建 HeartDance 客户端
	client := pb.NewMasterServerClient(conn)
	context, cancle := context.WithCancel(context.Background())
	defer cancle()
	// 模拟的节点ID
	nodeID := Args()

	// 启动心跳发送
	for {
		r, err := client.HeartDance(context, &pb.Signal{Mechine: nodeID})
		if err != nil {
			log.Println(err)
			break
		}
		MasterOnline = (r.GetOniline() == 1)
		time.Sleep(5 * time.Second)
	}
	conn.Close()
}

func Args() string {
	args := os.Args
	if len(args) != 2 {
		log.Println("Usage: ./node <nodeID>")
		os.Exit(1)
	}
	return args[1]
}

func FileStorageServer() {
	lis, err := net.Listen("tcp", ":45679")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServerServer(s, &NodeStorageServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("server can't run,err:%s", err)
	}
}
