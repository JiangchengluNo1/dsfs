package node

import (
	"context"
	"fmt"
	"log"
	"time"

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
	client, err = grpc.NewClient("localhost:56789", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	c = noding.NewNodingClient(client)
	ctx = context.Background()
	res, err := c.Wake(ctx, &noding.WakeUp{Files: files})
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Success)
}

func Healthing(ip string) {
	defer client.Close()
	for {
		res, err := c.Heart(ctx, &noding.Hearting{Ip: ip})
		if err != nil {
			log.Printf("Heart error: %v\n", err)
			break
		}
		fmt.Println(res.Success)
		time.Sleep(5 * time.Second)
	}
}
