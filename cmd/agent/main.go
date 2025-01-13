package main

import (
	"log"
	"net"

	"github.com/mfreeman451/homemon"
	"github.com/mfreeman451/homemon/pkg/agent"
	pb "github.com/mfreeman451/homemon/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	checkers := map[string]homemon.Checker{
		"nginx": &agent.ProcessChecker{ProcessName: "nginx"},
		"ssh":   &agent.PortChecker{Host: "localhost", Port: 22},
	}

	server := grpc.NewServer()
	pb.RegisterAgentServiceServer(server, agent.NewServer(checkers))

	log.Printf("Agent server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
