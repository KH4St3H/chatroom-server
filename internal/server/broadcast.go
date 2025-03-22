package server

import (
	"github.com/kh4st3h/chatroom-server/internal/server/types/request"
)

func (s *Server) Broadcast(u request.AuthenticatedUserRequest) {
	for username, _ := range s.Connections {
		if username == u.Username {
			continue
		}
		s.BroadcastTo(username, u.Message)
	}
}

func (s *Server) BroadcastTo(username string, msg string) {
	conn, found := s.Connections[username]
	if !found {
		logger.Errorw("user does not exist", "username", username)
		return
	}
	err := conn.EncryptedWrite([]byte(msg))
	if err != nil {
		logger.Error("failed to send data to user")
		conn.GoOffline()
		s.Leave(conn)
	}
}
