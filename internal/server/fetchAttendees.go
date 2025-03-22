package server

import (
	"github.com/kh4st3h/chatroom-server/internal/server/types/connection"
	"strings"
)

func (s *Server) fetchAttendees(conn connection.Conn) {
	usernames := make([]string, 0, len(s.Connections))

	for username, _ := range s.Connections {
		usernames = append(usernames, username)
	}

	msg := "Here are the list of attendees:\n" + strings.Join(usernames, ", ")
	s.BroadcastTo(&conn, msg)
}
