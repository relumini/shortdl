package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/relumini/shortdl/database"
	pb "github.com/relumini/shortdl/protos"
	"github.com/relumini/shortdl/routes"
	"google.golang.org/grpc"
)

const (
	serverAddress = "localhost:50051"
)

var Client pb.DownloadShortClient

func main() {
	_, err := database.ConnectDB()
	if err != nil {
		fmt.Print(err)
	}
	conn, err := grpc.NewClient(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close() // Close the connection on exit

	// Create a gRPC client for the DownloadShort service
	Client = pb.NewDownloadShortClient(conn)
	// Auto Migrate the model
	r := gin.Default()
	routes.InitRoute(r, Client)

	r.Run("localhost:8080")
}
