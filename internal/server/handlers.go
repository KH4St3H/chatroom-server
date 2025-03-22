package server

import (
	"bytes"
	"context"
	"github.com/kh4st3h/chatroom-server/internal/constants"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"github.com/kh4st3h/chatroom-server/internal/server/handler"
	"github.com/kh4st3h/chatroom-server/internal/server/types/connection"
	"github.com/kh4st3h/chatroom-server/internal/server/types/request"
	"go.uber.org/zap"
	"net"
)

var logger *zap.SugaredLogger

func init() {
	logger = log.NewLogger().Sugar()
}

func (s *Server) HandleUserMessages(conn *connection.Conn) {
	for {
		data, err := conn.ReadAndDecrypt()
		if err != nil {
			conn.GoOffline()
			logger.Errorf("Failed to read data from user: %v", err)
			return
		}
		if data == "" {
			continue
		}
		userMessage := request.New(conn.GetUsername(), data)
		s.UserRequests <- userMessage
	}
}

func (s *Server) HandleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Fatal connection received")
		}
	}()
	packet := make([]byte, 4096)
	count, err := conn.Read(packet)
	packet = packet[:count]
	if err != nil {
		logger.Errorf("Error reading from connection: %v", err)
		conn.Close()
		return
	}
	logger.Debugw("first packet received", "packet", string(packet))
	ctx = context.WithValue(ctx, "data", packet)
	newConn := connection.New(conn)

	if bytes.HasPrefix(packet, []byte(constants.REGISTRATION_MSG)) {
		handler.HandleRegistration(*newConn, ctx)
		conn.Close()
	} else if bytes.HasPrefix(packet, []byte(constants.LOGIN_MSG)) {
		success := handler.HandleLogin(newConn, ctx)
		if !success {
			conn.Close()
			return
		}
		go s.HandleUserMessages(newConn)
		s.Join(newConn)
	}
	return
}

func (s *Server) HandleNewConnections(ctx context.Context, connections chan *net.Conn) {
	for conn := range connections {
		logger.Infow("Incoming connection", "connection", conn)
		go s.HandleConnection(ctx, *conn)
	}
}

func (s *Server) HandleRequests() {
	for {
		req := <-s.UserRequests
		switch req.Message {
		case constants.REGISTRATION_MSG:

		}

	}

}
