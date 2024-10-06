package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	Model "github.com/mahaonan001/dsfs/cmd/master/model"
	pb "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
)

// maxWaiting 最大等待时间
const maxWaiting = 10 * time.Second
const mutexMaxWaitingTime = 2 * time.Second

type Heart struct {
	NodeOnline map[string]time.Time
	sync.Mutex
}

type Master struct {
	Heart
	Model.StreamPoint
	Model.FsMap
	pb.UnimplementedMasterServerServer
}

// GetFile implements filetransfer.MasterServerServer.
// Subtle: this method shadows the method (UnimplementedMasterServerServer).GetFile of Master.UnimplementedMasterServerServer.
func (m *Master) GetFile(context.Context, *pb.File2Get) (*pb.FileGetRes, error) {
	panic("unimplemented")
}

// HeartDance implements filetransfer.MasterServerServer.
func (m *Master) HeartDance(context.Context, *pb.Signal) (*pb.Alive, error) {
	panic("unimplemented")
}

// MasterWakeUp implements filetransfer.MasterServerServer.
// Subtle: this method shadows the method (UnimplementedMasterServerServer).MasterWakeUp of Master.UnimplementedMasterServerServer.
func (m *Master) MasterWakeUp(context.Context, *pb.MWU) (*pb.Alive, error) {
	panic("unimplemented")
}

// UploadFile implements filetransfer.MasterServerServer.
// Subtle: this method shadows the method (UnimplementedMasterServerServer).UploadFile of Master.UnimplementedMasterServerServer.
func (m *Master) UploadFile(context.Context, *pb.File2Up) (*pb.FileUpRes, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedMasterServerServer implements filetransfer.MasterServerServer.
// Subtle: this method shadows the method (UnimplementedMasterServerServer).mustEmbedUnimplementedMasterServerServer of Master.UnimplementedMasterServerServer.
// func (m *Master) mustEmbedUnimplementedMasterServerServer() {
// 	panic("unimplemented")
// }

var master = &Master{Heart: Heart{NodeOnline: make(map[string]time.Time)}}

// ChangeNodeMessage 更改node的信息，超时返回错误
func (h *Heart) ChangeNodeMessage(id string) error {
	st := time.Now()
	for {
		switch {
		case h.TryLock():
			{
				fl := len(h.NodeOnline)
				h.NodeOnline[id] = time.Now()
				if fl == len(h.NodeOnline)-1 {
					log.Printf("%s is online......", id)
				}
				h.Unlock()
				return nil
			}
		case time.Since(st) > mutexMaxWaitingTime:
			{
				return errors.New("get mutex time out,please check the master")
			}
		}
	}
}

// NodeComeOn 当一个node连接到master时，master会调用这个函数
func (h *Heart) HeartDance(ctx context.Context, req *pb.Signal) (*pb.Alive, error) {
	err := master.ChangeNodeMessage(req.Mechine)
	if err != nil {
		return nil, err
	}
	return &pb.Alive{Oniline: 1}, nil
}

// CheckNodeStatus 检查node的状态
func (h *Heart) CheckNodeStatus() {
	for {
		time.Sleep(5 * time.Second)
		h.Lock()
		for k := range h.NodeOnline {
			t := time.Now()
			if t.Sub(h.NodeOnline[k]) > maxWaiting {
				delete(h.NodeOnline, k)
				log.Printf("node %s is offline", k)
			}
		}
		h.Unlock()
	}
}

// func ()

func main() {
	go master.CheckNodeStatus()
	lis, err := net.Listen("tcp", ":45678")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMasterServerServer(s, &Master{})
	fmt.Println("master is running at :45678")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
