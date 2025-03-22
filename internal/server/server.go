package server

import (
	"context"
	"fmt"
	"github.com/kh4st3h/chatroom-server/internal/config"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"github.com/kh4st3h/chatroom-server/internal/server/types/request"
	"net"
)

type Server struct {
	Cfg          *config.Config
	UserRequests chan *request.AuthenticatedUserRequest
}

func NewServer(cfg *config.Config) *Server {
	return &Server{cfg, make(chan *request.AuthenticatedUserRequest)}
}

func (s *Server) Run() error {
	logger := log.NewLogger().Sugar()
	fullListenAddr := fmt.Sprintf("%s:%d", s.Cfg.ListenAddr, s.Cfg.ListenPort)
	logger.Infow("Starting server", "listenAddr", fullListenAddr)
	l, err := net.Listen("tcp", fullListenAddr)
	if err != nil {
		return err
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			logger.Warnw("failed to close listener", "err", err)
		}
	}(l)
	connections := make(chan *net.Conn)
	ctx := context.Background()
	go s.HandleNewConnections(ctx, connections)
	go s.HandleRequests()
	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Errorf("Error accepting incoming connection: %v", err)
		}
		if err != nil {
			logger.Errorf("Error setting deadline: %v", err)
			return err
		}
		connections <- &conn
	}
}
