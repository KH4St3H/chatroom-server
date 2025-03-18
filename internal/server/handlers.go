package server

import (
	"bytes"
	"github.com/kh4st3h/chatroom-server/internal/constants"
	"github.com/kh4st3h/chatroom-server/internal/server/handlers"
	"net"
)

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()
	packet := make([]byte, 4096)
	count, err := conn.Read(packet)
	packet = packet[:count]
	if err != nil {
		s.Logger.Errorf("Error reading from connection: %v", err)
		return
	}
	s.Logger.Debugw("first packet received", "packet", string(packet))
	if bytes.HasPrefix(packet, []byte(constants.REGISTRATION_MSG)) {
		err := handlers.HandleRegistration(s.Context, packet)
		if err != nil {
			_, err := conn.Write([]byte(err.Error()))
			if err != nil {
				s.Logger.Errorf("Error writing to connection: %v", err)
				return
			}
			return
		}
		s.Logger.Info("User created")
		_, err = conn.Write([]byte("User created successfully!"))
		if err != nil {
			s.Logger.Errorf("Error writing to connection: %v", err)
		}
	}
	return
}

func (s *Server) HandleNewConnections(connections chan *net.Conn) {
	for conn := range connections {
		s.Logger.Infow("Incoming connection", "connection", conn)
		go s.HandleConnection(*conn)
	}
}
