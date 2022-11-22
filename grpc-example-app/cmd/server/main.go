package main

import (
	"log"
	"net"

	"github.com/keploy/go-sdk/integrations/kgrpcserver"
	"github.com/keploy/go-sdk/keploy"
	"github.com/koddr/example-go-grpc-server/pkg/adder"
	api "github.com/koddr/example-go-grpc-server/pkg/api/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Create new gRPC server instance
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "grpc-nested-app",
			Port: "8009",
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:6789/api",
		},
	})
	// s := grpc.NewServer()
	srv := &adder.GRPCServer{}
	s := grpc.NewServer(kgrpcserver.UnaryInterceptor(k))
	reflection.Register(s)
	// Register gRPC server
	api.RegisterAdderServer(s, srv)

	// Listen on port 8080
	l, err := net.Listen("tcp", ":8009")
	if err != nil {
		log.Fatal(err)
	}

	// Start gRPC server
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
