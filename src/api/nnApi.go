package api

import (
	"context"
	"fmt"
	"proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type NNOut struct {
	value float32
	p     []float32
}

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func nnTrain(trainFromIndex int32) int32 {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	req := &proto.TrainFromIndex{Uid: trainFromIndex}

	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.Train(ctx, req)
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	conn.Close()

	return response.Status
}

func nnForwardPass(game *Connect4) NNOut {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	board := game.GetBoardFlatInt32()
	req := &proto.BoardRequest{Board: board}
	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.ForwardPass(ctx, req)
	conn.Close()

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

/*
  Train the NN using the play data from the last iteration
*/
func TrainFromLastIteration(uid int32) {
	nnTrain(uid)
}
