package server

import (
	"bytes"
	"context"
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/constants"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"github.com/kh4st3h/chatroom-server/internal/server/handler"
	"go.uber.org/zap"
	"net"
	"time"
)

type Conn struct {
	conn net.Conn
	ctx  context.Context
}

var logger *zap.SugaredLogger

func init() {
	logger = log.NewLogger().Sugar()

}

func (c *Conn) Write(data []byte) (int, error) {
	count, err := c.conn.Write(data)
	if err != nil {
		logger.Errorw("Failed to write data to client", "data", data)
	}
	return count, err
}

func (s *Server) HandleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	packet := make([]byte, 4096)
	count, err := conn.Read(packet)
	packet = packet[:count]
	if err != nil {
		logger.Errorf("Error reading from connection: %v", err)
		return
	}
	logger.Debugw("first packet received", "packet", string(packet))
	ctx = context.WithValue(ctx, "data", packet)
	newConn := &Conn{conn: conn, ctx: ctx}

	if bytes.HasPrefix(packet, []byte(constants.REGISTRATION_MSG)) {
		handler.HandleRegistration(newConn, ctx)
	} else if bytes.HasPrefix(packet, []byte(constants.LOGIN_MSG)) {
		success := handler.HandleLogin(newConn, ctx)
		if !success {
			return
		}
	}
	return
}

func (s *Server) HandleNewConnections(ctx context.Context, connections chan *net.Conn) {
	for conn := range connections {
		logger.Infow("Incoming connection", "connection", conn)
		ctx, cancel := context.WithTimeout(ctx, s.Cfg.ConnectionTimeout*time.Millisecond)
		defer cancel()
		go s.HandleConnection(ctx, *conn)
		go func() {
			select {
			case <-ctx.Done():
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					logger.Errorf("Timeout reached: %v", ctx.Err())
				} else {
					logger.Infow("Finished processing connection", "connection", conn)
				}
			}
		}()
	}
}
