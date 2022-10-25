package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"github.com/gin-gonic/gin"
	"github.com/keploy/go-sdk/integrations/kgin/v1"
	"github.com/keploy/go-sdk/integrations/kgrpc"
	"github.com/keploy/go-sdk/keploy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	address = "localhost:50051"
)

var conn *grpc.ClientConn

func getURL(c *gin.Context) {
	// fmt.Println("\n\n", conn, "\n ")
	client := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Hour)
	defer cancel()

	fmt.Println("\n\n", client, "\n ")
	res, err := client.CreateNewUser(ctx, &pb.NewUser{Name: "Ritik", Age: 21})
	if err != nil {
		log.Fatalf("not created: %v", err)
	}
	eRr, _ := status.FromError(err)
	log.Printf(`Print
	Name: %v
	Age: %v
	Id: %v
	Error type: %v`, res.GetName(), res.GetAge(), res.GetId(), eRr.String())
}

func main() {
	r := gin.New()
	port := "8001"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "grpc-client-app",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:6789/api",
		},
	})
	kgin.GinV1(k, r)
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure(), kgrpc.WithClientUnaryInterceptor(k))
	if err != nil {
		log.Fatalf("Did not connect : %v", err)
	}
	defer conn.Close()
	r.GET("/param", getURL)
	r.Run(":" + port)
}
