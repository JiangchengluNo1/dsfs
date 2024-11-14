package node

import (
	"context"

	noding "github.com/mahaonan001/dsfs/proto/healthing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Address = ":8520"
)

var (
	client *grpc.ClientConn
	c      noding.NodingClient
	ctx    context.Context
	err    error
)

func WakeUp(files []string) {
	client, err = grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	c = noding.NewNodingClient(client)
	ctx = context.Background()
	c.Wake(ctx, &noding.WakeUp{Files: files})
}

func Healthing(number int) {
	c.Heart(ctx, &noding.Hearting{BeatNumber: int32(number)})
}
