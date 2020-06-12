package api

import (
	"context"
	"fmt"
	"proto"
	"time"

	"google.golang.org/grpc"
)

type NNOut struct {
	value float32
	p     []float32
}

func nnForwardPass(game *Connect4) NNOut {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	board := game.GetBoardFlatInt32()
	req := &proto.BoardRequest{Board: board}
	deadlineMs := 10000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.ForwardPass(ctx, req)

	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	var nnOut NNOut
	//fmt.Printf("Response %v", response.Result)

	nnOut.value = float32(response.Result[0])
	for i := 1; i <= 7; i++ {
		//fmt.Println(i)
		nnOut.p = append(nnOut.p, float32(response.Result[i]))
	}

	return nnOut
}
