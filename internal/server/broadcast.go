package server

import (
	"github.com/kh4st3h/chatroom-server/internal/server/types/request"
)

func (s *Server) Broadcast(u request.AuthenticatedUserRequest) {
	for username, conn := range s.Connections {
		if username == u.Username {
			continue
		}

		err := conn.EncryptedWrite([]byte(u.Message))
		if err != nil {
			logger.Error("failed to send data to user")
			conn.GoOffline()
			s.Leave(conn)
		}
	}
}
