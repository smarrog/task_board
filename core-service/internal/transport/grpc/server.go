package grpc

import (
	"net"

	"github.com/rs/zerolog"
	v1 "github.com/smarrog/task-board/shared/proto/base/v1"
	"google.golang.org/grpc"
)

type Server struct {
	log    *zerolog.Logger
	server *grpc.Server
}

func NewServer(
	log *zerolog.Logger,
	boardsHandler *BoardsHandler,
	columnsHandler *ColumnsHandler,
	tasksHandler *TasksHandler,
) *Server {
	s := grpc.NewServer()

	RegisterHealth(s)

	v1.RegisterBoardsServiceServer(s, boardsHandler)
	v1.RegisterColumnsServiceServer(s, columnsHandler)
	v1.RegisterTasksServiceServer(s, tasksHandler)

	return &Server{
		log:    log,
		server: s,
	}
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.log.Info().Str("addr", addr).Msg("gRPC server started")
	return s.server.Serve(lis)
}

func (s *Server) Stop() {
	s.log.Info().Msg("stopping gRPC server")
	s.server.GracefulStop()
}
