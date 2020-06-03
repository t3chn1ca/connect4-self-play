package main

import (
	"api"
	"context"
	"fmt"
	"proto"
	"time"

	"google.golang.org/grpc"
)

//For generating boards for test, remove later
func setupGame(game *api.Connect4, moves []int) *api.Connect4 {

	for _, move := range moves {

		//fmt.Println("===============================")
		//fmt.Printf("Player to move %s\n", game.PlayerToString(game.GetPlayerToMove()))
		game.PlayMove(move)
		//fmt.Printf("Player who just Moved %s\n", game.PlayerToString(game.GetPlayerWhoJustMoved()))
	}
	//fmt.Printf("Board Index = %s", game.GetBoardIndex().String())
	return game

}

func main() {
	var game = api.NewConnect4()
	moves := []int{1, 0, 2, 4, 3, 0, 6, 1, 3, 2, 4, 6, 2, 0, 4, 1, 6, 3, 0, 4, 5, 3, 6, 3} //, 6, 2}
	game = setupGame(game, moves)

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	a := game.GetBoardFlatInt32()
	fmt.Println("a = %v", a)
	//var b int = 50
	req := &proto.BoardRequest{Board: a}
	deadlineMs := 10000
	//Increase timeout as the server on init takes few secs to return ( its loading nn into GPU memory)
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, _ := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.ForwardPass(ctx, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for _, i := range response.Result {
		fmt.Printf(" Response %f\n", float32(i))
	}

}
