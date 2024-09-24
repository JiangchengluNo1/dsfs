package main

import (
	"errors"
	"log"
	"sync"
	"time"

	pb "github.com/mahaonan001/dsfs/proto"
)

// maxWaiting 最大等待时间
const maxWaiting = 10 * time.Second

// Hearting 心跳服务
type Hearting struct {
	OnlineNode map[string]time.Time
	pb.UnimplementedHeartDanceServer
	sync.Mutex
}

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
func (h *Hearting) NodeComeOn(req *pb.Signal) (*pb.Alive, error) {
	return &pb.Alive{Oniline: 1}, h.ChangeNodeMessage(req.Mechine)
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

}
