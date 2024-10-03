package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var MasterOnline = false

func main() {
	// 设置与 Master 的 gRPC 连接
	conn, err := grpc.NewClient("localhost:45678", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()

	// 创建 HeartDance 客户端
	client := pb.NewHeartDanceClient(conn)
	context, cancle := context.WithCancel(context.Background())
	defer cancle()
	// 模拟的节点ID
	nodeID := Args()

	// 启动心跳发送
	for {
		r, err := client.HeartDance(context, &pb.Signal{Mechine: nodeID})
		if err != nil {
			log.Println(err)
		}
		MasterOnline = (r.GetOniline() == 1)
		time.Sleep(5 * time.Second)
	}
}

func Args() string {
	args := os.Args
	if len(args) != 2 {
		log.Println("Usage: ./node <nodeID>")
		os.Exit(1)
	}
	return args[1]
}
