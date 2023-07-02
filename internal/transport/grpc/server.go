package grpc

import (
	pb_static "github.com/ensiouel/basket-contract/gen/go/static/v1"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	grpcServer *grpc.Server
	logger     *slog.Logger
}

func New(logger *slog.Logger) *Server {
	grpcServer := grpc.NewServer()

	return &Server{
		grpcServer: grpcServer,
		logger:     logger,
	}
}

func (server *Server) Register(staticServer pb_static.StaticServer) *Server {
	pb_static.RegisterStaticServer(server.grpcServer, staticServer)

	return server
}

func (server *Server) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return server.grpcServer.Serve(listener)
}

func (server *Server) Stop() {
	server.grpcServer.GracefulStop()
}
