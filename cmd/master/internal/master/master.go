package master

import (
	"context"

	noding "github.com/mahaonan001/dsfs/proto/healthing"
)

type MasterServer struct {
	noding.UnimplementedNodingServer
}

// //////////////////////////////////////
var file2number map[string][]int32

// //////////////////////////////////////
func (m *MasterServer) HealthCheck(ctx context.Context, in *noding.Hearting) (*noding.HeartingResponse, error) {
	beatnumber := in.GetBeatNumber()
	return &noding.HeartingResponse{Success: CheckNodeOnline(beatnumber)}, nil
}

func (m *MasterServer) WakeUp(ctx context.Context, in *noding.WakeUp) (*noding.WakeUpResponse, error) {
	files := in.GetFiles()
	number := in.GetNumber()
	for _, v := range files {
		file2number[v] = append(file2number[v], number)
	}
	return &noding.WakeUpResponse{Success: true}, nil
}
func CheckNodeOnline(beatnumber int32) bool {
	return true
}
