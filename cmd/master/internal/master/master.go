package master

import (
	"context"
	"fmt"
	"sync"
	"time"

	noding "github.com/mahaonan001/dsfs/proto/healthing"
)

type MasterServer struct {
	noding.UnimplementedNodingServer
	NodeClient map[int32]time.Time // client onilne or not
	FileOnNode map[string]int32    // file to client
	sync.RWMutex
}

// //////////////////////////////////////
var file2number map[string][]int32
var Master2CF MasterServer

// //////////////////////////////////////
func (m *MasterServer) Heart(ctx context.Context, in *noding.Hearting) (*noding.HeartingResponse, error) {
	beatnumber := in.GetBeatNumber()
	if !m.CheckNodeOnline(beatnumber) {
		fmt.Println("Node", beatnumber, "is online")
	}
	m.NodeClient[beatnumber] = time.Now()
	return &noding.HeartingResponse{Success: m.CheckNodeOnline(beatnumber)}, nil
}

func (m *MasterServer) Wake(ctx context.Context, in *noding.WakeUp) (*noding.WakeUpResponse, error) {
	files := in.GetFiles()
	number := in.GetNumber()
	for _, v := range files {
		file2number[v] = append(file2number[v], number)
	}
	return &noding.WakeUpResponse{Success: true}, nil
}
func (m *MasterServer) CheckNodeOnline(beatnumber int32) bool {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.NodeClient[beatnumber]
	return ok
}

func ClientOffSound() {

	for {
		for k := range Master2CF.NodeClient {
			Master2CF.Lock()
			if time.Since(Master2CF.NodeClient[k]) > time.Second*10 {
				red := "\033[31m"
				reset := "\033[0m"
				fmt.Printf("%sNode %v is out of line\n%s", red, k, reset)
				delete(Master2CF.NodeClient, k)
			}
			Master2CF.Unlock()
		}
		time.Sleep(time.Second * 5)
	}
}

func init() {
	Master2CF = MasterServer{
		NodeClient: make(map[int32]time.Time),
		FileOnNode: make(map[string]int32),
	}
	file2number = make(map[string][]int32)
}
