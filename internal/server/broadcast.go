package server

import (
	"github.com/kh4st3h/chatroom-server/internal/server/types/connection"
	"github.com/kh4st3h/chatroom-server/internal/server/types/request"
)

func (s *Server) Broadcast(u request.AuthenticatedUserRequest) {
	for username, conn := range s.Connections {
		if username == u.Username {
			continue
		}
		s.BroadcastTo(conn, u.Message)
	}
}

func (s *Server) BroadcastTo(conn *connection.Conn, msg string) {
	err := conn.EncryptedWrite([]byte(msg))
	if err != nil {
		logger.Error("failed to send data to user")
		conn.GoOffline()
		s.Leave(conn)
	}
}
