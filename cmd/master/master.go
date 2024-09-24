package main

import (
	"errors"
	"log"
	"sync"
	"time"

	pb "github.com/mahaonan001/dsfs/proto"
)

const maxWaiting = 10 * time.Second

type Hearting struct {
	OnlineNode map[string]time.Time
	pb.UnimplementedHeartDanceServer
	sync.Mutex
}

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

// 实现文件上传
func (h *Hearting) NodeComeOn(req *pb.Signal) (*pb.Alive, error) {
	h.Lock()
	h.ChangeNodeMessage(req.Mechine)
	h.Unlock()
	return &pb.Alive{Oniline: 1}, nil
}

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

}
