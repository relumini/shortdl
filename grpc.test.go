package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	// Import the generated gRPC client code from your .proto file
	pb "github.com/relumini/shortdl/protos" // Replace with the actual path
)

const (
	// Server address (replace with your server's address and port)
	serverAddr = "localhost:50051"
)

var dc pb.DownloadShortClient

func initDc() {
	conn, err := grpc.NewClient(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close() // Close the connection on exit

	// Create a gRPC client for the DownloadShort service
	dc = pb.NewDownloadShortClient(conn)
}

func test() {
	initDc()
	// Create a connection to the server
	// conn, err := grpc.NewClient(serverAddress, grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalf("Failed to dial server: %v", err)
	// }
	// defer conn.Close() // Close the connection on exit

	// // Create a gRPC client for the DownloadShort service
	// client := pb.NewDownloadShortClient(conn)

	// Example call to DownTiktok
	url := "https://www.tiktok.com/@pt.memeindo/video/7360948743601999109" // Replace with the actual TikTok URL
	request := &pb.ParamsRequest{
		Url: url,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set a timeout for the request
	defer cancel()

	response, err := dc.DownTiktok(ctx, request)
	if err != nil {
		log.Fatalf("Failed to call DownTiktok: %v", err)
	}

	fmt.Println("DownloadTiktok response:", response.Status)

	// You can similarly call DownYoutube with appropriate URL and request message

}
