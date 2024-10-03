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

type Hearting struct {
	OnlineNode map[string]time.Time
	pb.UnimplementedMasterServerServer
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
				fl := len(h.OnlineNode)
				h.OnlineNode[id] = time.Now()
				if fl == len(h.OnlineNode)-1 {
					log.Printf("%s is online......", id)
				}
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
func (h *Hearting) HeartDance(ctx context.Context, req *pb.Signal) (*pb.Alive, error) {
	// log.Printf("node %s is online", req.Mechine)
	err := master.ChangeNodeMessage(req.Mechine)
	if err != nil {
		return nil, err
	}
	return &pb.Alive{Oniline: 1}, nil
}

func (h *Hearting) MasterWakeUp(ctx context.Context, req *pb.MWU) (*pb.Alive, error) {

	if err := master.ChangeNodeMessage(req.Sig.Mechine); err != nil {
		return nil, err
	}
	return &pb.Alive{Oniline: 1}, nil
}

// CheckNodeStatus 检查node的状态
func (h *Hearting) CheckNodeStatus() {
	for {
		time.Sleep(5 * time.Second)
		h.Lock()
		for k := range h.OnlineNode {
			t := time.Now()
			if t.Sub(h.OnlineNode[k]) > 10*time.Second {
				delete(h.OnlineNode, k)
				log.Printf("node %s is offline", k)
			}
		}
		h.Unlock()
	}
}

func main() {
	go master.CheckNodeStatus()
	lis, err := net.Listen("tcp", ":45678")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterServerServer(s, &Hearting{})
	fmt.Println("master is running at :45678")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
