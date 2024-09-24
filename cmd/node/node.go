package main

import (
	"context"
	"log"

	pb "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 设置与 Master 的 gRPC 连接
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()

	// 创建 HeartDance 客户端
	client := pb.NewHeartDanceClient(conn)
	context, cancle := context.WithCancel(context.Background())
	defer cancle()
	// 模拟的节点ID
	nodeID := "node-1"

	// 启动心跳发送
	r, err := client.HeartDance(context, &pb.Signal{Mechine: nodeID})
	if err != nil {
		log.Println(err)
	}
	log.Println(r.GetOniline())
}
