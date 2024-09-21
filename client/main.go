package main

import (
	"context"
	"fmt"
	"time"

	solider "github.com/mahaonan001/dsfs/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Gctx context.Context
var Gcancle context.CancelFunc
var Gclient solider.EatHandlerClient

func Clinet2Server() (*grpc.ClientConn, error) {
	client, err := grpc.NewClient("localhost:9876", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	Gclient = solider.NewEatHandlerClient(client)
	return client, nil
}

func main() {
	cli, err := Clinet2Server()
	defer func() {
		err := cli.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic("failed to create new grpc client")
	}
	Gctx, Gcancle = context.WithTimeout(context.Background(), time.Hour)
	defer Gcancle()
	c := solider.NewEatHandlerClient(cli)
	r, err := c.FileHandler(Gctx, &solider.InFile{Datas: []byte{12, 31, 45, 2}, Unionid: "mahaonan"})
	if err != nil {
		panic("failed to get response,err: " + err.Error())
	}
	fmt.Println(r.GetCode(), r.GetMachine(), r.GetNumberMachine())
}
