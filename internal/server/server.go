package server

import (
	"fmt"
	"github.com/kh4st3h/chatroom-server/internal/config"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"go.uber.org/zap"
	"net"
)

type Server struct {
	Logger      *zap.SugaredLogger
	Config      *config.Config
	DBManager   *db.Manager
	Connections chan *net.Conn
}

func NewServer(cfg *config.Config, logger *zap.SugaredLogger, dbManager *db.Manager) *Server {
	connections := make(chan *net.Conn)
	return &Server{
		Logger:      logger,
		Config:      cfg,
		Connections: connections,
		DBManager:   dbManager,
	}
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
	go s.HandleNewConnections()
	for {
		conn, err := l.Accept()
		if err != nil {
			s.Logger.Errorf("Error accepting incoming connection: %v", err)
		}
		s.Connections <- &conn
	}
}
