package main

import (
	"net"

	"github.com/mahaonan001/dsfs/cmd/master/internal/master"
	noding "github.com/mahaonan001/dsfs/proto/healthing"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":56789")
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	go master.ClientOffSound()
	noding.RegisterNodingServer(s, &master.Master2CF)
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
