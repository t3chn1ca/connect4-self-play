package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"proto"
	"shared"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func NnLoadBestModelToGpu() int32 {
	fmt.Println("NnLoadBestModelToGpu: Loading best model to GPU ")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)

	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	args := &proto.NoArg{}
	response, err := client.LoadBestModelToGpu(ctx, args)
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	conn.Close()

	return response.Status
}
func NnSaveTrainedModelToBest() int32 {
	fmt.Println("NnSaveTrainedModelToBest: Saving trained model to best")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)

	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	args := &proto.NoArg{}
	response, err := client.SaveCurrentModelToBest(ctx, args)
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	conn.Close()

	return response.Status
}

func NnStopCpuModel() int32 {
	fmt.Println("NnStopCpuModel: Stopping model to cpu ")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)

	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	args := &proto.NoArg{}
	response, err := client.StopBestModelCpu(ctx, args)
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	conn.Close()

	return response.Status
}

func NnLoadCpuModel() int32 {
	fmt.Println("NnLoadCpuModel: Loading model to cpu ")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)

	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	args := &proto.NoArg{}
	response, err := client.LoadBestModelToCpu(ctx, args)
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	conn.Close()

	return response.Status
}

func NnTrain(trainFromIndex int32, trainToIndex int32) int32 {
	fmt.Println("nnTrain: Training model on GPU")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	req := &proto.TrainFromIndex{UidFrom: trainFromIndex, UidTo: trainToIndex}

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

func NnForwardPassMultiBoard(game *Connect4, port int) ([MAX_CHILD_NODES + 1]shared.NNOut, [MAX_CHILD_NODES + 1]big.Int) {
	portStr := "localhost:" + strconv.Itoa(port)
	conn, err := grpc.Dial(portStr, grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	boardsJson, boardIndexes := game.GetValidFlatBoardsFromPosition()
	//fmt.Printf("%s\n", string(boardsJson))
	//fmt.Printf("%v\n", boardIndexes)
	req := &proto.MultiBoardRequest{Json: string(boardsJson)}
	deadlineMs := 1000000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.ForwardPassMultiBoard(ctx, req)
	conn.Close()
	if err != nil {
		fmt.Println("ERROR!: Connection to server failed, check server is started!")
		panic(err)
	}
	//                        7 child + 1 Parent    1v + 7 p's
	var multiBoardNNResponse [MAX_CHILD_NODES + 1][MAX_CHILD_NODES + 1]float32

	err = json.Unmarshal([]byte(response.Json), &multiBoardNNResponse)

	if err != nil {
		fmt.Println("ERROR!: Decoding Json for multiboard response")
		panic(err)
	}

	var nnOutArray [MAX_CHILD_NODES + 1]shared.NNOut
	//boardIndex is used to index the 8 boards which are returned by NN
	for boardIndex := 0; boardIndex < MAX_CHILD_NODES+1; boardIndex++ {
		nnOutArray[boardIndex].Value = float32(multiBoardNNResponse[boardIndex][0])
		//actionIndex starts at 1-7, 0th index is used for value
		for actionIndex := 1; actionIndex <= MAX_CHILD_NODES; actionIndex++ {
			nnOutArray[boardIndex].P = append(nnOutArray[boardIndex].P, multiBoardNNResponse[boardIndex][actionIndex])
		}
	}

	return nnOutArray, boardIndexes
}

func NnForwardPass(game *Connect4, port int) shared.NNOut {
	portStr := "localhost:" + strconv.Itoa(port)
	conn, err := grpc.Dial(portStr, grpc.WithInsecure()) //, grpc.WithKeepaliveParams(kacp))
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

	var nnOut shared.NNOut
	//fmt.Printf("Response %v", response.Result)

	nnOut.Value = float32(response.Result[0])
	for i := 1; i <= MAX_CHILD_NODES; i++ {
		//fmt.Println(i)
		nnOut.P = append(nnOut.P, float32(response.Result[i]))
	}

	return nnOut
}

/*
  Train the NN using the play data from the last iteration
*/
func TrainFromLastIteration(uidFrom int32, uidTo int32) {
	NnTrain(uidFrom, uidTo)
}
