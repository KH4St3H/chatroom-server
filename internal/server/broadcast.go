package server

func (s *Server) Broadcast(sender string, message string) {
	for username, _ := range s.Connections {
		if username == sender {
			continue
		}
		s.BroadcastTo(username, message)
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
		s.ConnectionFail(conn)
	}
}
