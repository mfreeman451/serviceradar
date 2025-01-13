package agent

import (
	"context"

	"github.com/mfreeman451/homemon/pkg/checker"
	"github.com/mfreeman451/homemon/proto"
)

type Server struct {
	pb.UnimplementedAgentServiceServer
	checkers map[string]checker.Checker
}

func NewServer(checkers map[string]checker.Checker) *Server {
	return &Server{
		checkers: checkers,
	}
}

func (s *Server) GetStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	c, exists := s.checkers[req.ServiceName]
	if !exists {
		return &proto.StatusResponse{
			Available: false,
			Message:   "Service not found",
		}, nil
	}

	available, msg := c.Check(ctx)
	return &proto.StatusResponse{
		Available: available,
		Message:   msg,
	}, nil
}
