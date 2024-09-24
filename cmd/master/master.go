package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
)

// maxWaiting 最大等待时间
const maxWaiting = 10 * time.Second

type server struct {
	pb.UnimplementedHeartDanceServer
}

type Hearting struct {
	OnlineNode map[string]time.Time
	sync.Mutex
}

var master = &Hearting{OnlineNode: make(map[string]time.Time)}

// ChangeNodeMessage 更改node的信息，超时返回错误
func (h *Hearting) ChangeNodeMessage(id string) error {
	st := time.Now()
	for {
		switch {
		case h.TryLock():
			{
				h.OnlineNode[id] = st
				h.Unlock()
				return nil
			}
		case time.Since(st) > maxWaiting:
			{
				return errors.New("get mutex time out,please check the master")
			}
		}
	}
}

// NodeComeOn 当一个node连接到master时，master会调用这个函数
func (s *server) HeartDance(ctx context.Context, req *pb.Signal) (*pb.Alive, error) {
	log.Printf("node %s is online", req.Mechine)
	return &pb.Alive{Oniline: 1}, master.ChangeNodeMessage(req.Mechine)
}

func (s *server) MasterWakeUp(ctx context.Context, req *pb.MWU) (*pb.Alive, error) {
	log.Printf("master %s is online", req.Sig.Mechine)
	return &pb.Alive{Oniline: 1}, master.ChangeNodeMessage(req.Sig.Mechine)
}

// CheckNodeStatus 检查node的状态
func (h *Hearting) CheckNodeStatus() {
	for {
		time.Sleep(5 * time.Second)
		t := time.Now()
		for node, lastTime := range h.OnlineNode {
			switch {
			case h.TryLock():
				if time.Since(lastTime) > maxWaiting {
					delete(h.OnlineNode, node)
				}
				log.Printf("node %s is out of connection\f", node)
				h.Unlock()
			case time.Since(t) > maxWaiting:
				log.Fatal("get mutex time out,please check the master")
			}
		}
	}
}

func main() {
	go master.CheckNodeStatus()
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Listening on :8080")
	s := grpc.NewServer()
	pb.RegisterHeartDanceServer(s, &server{})
	log.Println("gRPC server registered")
	fmt.Println("master is running at :8080....")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
