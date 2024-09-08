package main

import (
	"context"
	"log"
	"net"

	solider "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
)

type server struct {
	solider.UnimplementedEatHandlerServer
}

func (s *server) FileHandler(ctx context.Context, in *solider.InFile) (*solider.Response, error) {
	log.Printf("Recive: %s", in.GetUnionid())
	return &solider.Response{Code: 200, Machine: []int32{1, 2}, NumberMachine: []int32{1, 2}}, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:9876")
	if err != nil {
		log.Fatal("tcp net start error: " + err.Error())
	}
	defer lis.Close()
	s := grpc.NewServer()
	solider.RegisterEatHandlerServer(s, &server{})
	s.Serve(lis)
}
