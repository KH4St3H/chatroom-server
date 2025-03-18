package server

import (
	"fmt"
	"github.com/kh4st3h/chatroom-server/internal/context"
	"net"
)

type Server struct {
	*context.Context
}

func NewServer(ctx *context.Context) *Server {
	return &Server{ctx}
}

func (s *Server) Run() error {
	fullListenAddr := fmt.Sprintf("%s:%d", s.Config.ListenAddr, s.Config.ListenPort)
	s.Logger.Infow("Starting server", "listenAddr", fullListenAddr)
	l, err := net.Listen("tcp", fullListenAddr)
	if err != nil {
		return err
	}
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			s.Logger.Warnw("failed to close listener", "err", err)
		}
	}(l)
	connections := make(chan *net.Conn)
	go s.HandleNewConnections(connections)
	for {
		conn, err := l.Accept()
		if err != nil {
			s.Logger.Errorf("Error accepting incoming connection: %v", err)
		}
		connections <- &conn
	}
}
