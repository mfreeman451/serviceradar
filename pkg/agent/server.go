package agent

import (
	"context"

	"github.com/mfreeman451/homemon"
	pb "github.com/mfreeman451/homemon/proto"
)

type Server struct {
	pb.UnimplementedAgentServiceServer
	checkers map[string]homemon.Checker
}

func NewServer(checkers map[string]homemon.Checker) *Server {
	return &Server{
		checkers: checkers,
	}
}

func (s *Server) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	checker, exists := s.checkers[req.ServiceName]
	if !exists {
		return &pb.StatusResponse{
			Available: false,
			Message:   "Service not found",
		}, nil
	}

	available, msg := checker.Check(ctx)
	return &pb.StatusResponse{
		Available: available,
		Message:   msg,
	}, nil
}
