package main

import (
	"context"
	"fmt"
	"proto"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)
	a := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	//var b int = 50
	req := &proto.BoardRequest{Board: a}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.ForwardPass(ctx, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf(" Response %d", int(response.Result))

	/*

		g := gin.Default()
		g.GET("/add/:a/:b", func(ctx *gin.Context) {
			a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
				return
			}
			b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
				return
			}
			req := &proto.Request{A: int64(a), B: int64(b)}
			if response, err := client.Add(ctx, req); err == nil {
				ctx.JSON(http.StatusOK, gin.H{
					"result": fmt.Sprint(response.Result),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
		})
		/*
			g.GET("/subtract/:a/:b", func(ctx *gin.Context) {
				a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
					return
				}
				b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
				if err != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
					return
				}
				req := &proto.Request{A: int64(a), B: int64(b)}
				if response, err := client.Subtract(ctx, req); err == nil {
					ctx.JSON(http.StatusOK, gin.H{
						"result": fmt.Sprint(response.Result),
					})
				} else {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"error": err.Error(),
					})
				}
			})
		if err := g.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}*/
}
