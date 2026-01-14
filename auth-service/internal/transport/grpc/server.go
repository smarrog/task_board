package grpc

import (
	"net"

	"github.com/rs/zerolog"
	authv1 "github.com/smarrog/task-board/shared/proto/auth/v1"
	"google.golang.org/grpc"
)

type Server struct {
	log *zerolog.Logger
	srv *grpc.Server
	h   *AuthHandler
}

func NewServer(log *zerolog.Logger, h *AuthHandler) *Server {
	gs := grpc.NewServer()
	authv1.RegisterAuthServiceServer(gs, h)
	return &Server{log: log, srv: gs, h: h}
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.log.Info().Str("addr", addr).Msg("gRPC server started")
	return s.srv.Serve(lis)
}

func (s *Server) Stop() {
	s.log.Info().Msg("stopping gRPC server")
	s.srv.GracefulStop()
}
