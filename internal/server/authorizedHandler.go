package server

import (
	"github.com/kh4st3h/chatroom-server/internal/actions/server"
	"github.com/kh4st3h/chatroom-server/internal/constants"
	"github.com/kh4st3h/chatroom-server/internal/server/types/connection"
	"github.com/kh4st3h/chatroom-server/internal/server/types/request"
)

func (s *Server) HandleUserMessages(conn *connection.Conn) {
	dataChan := make(chan string)
	for {
		go conn.ReadAndDecrypt(dataChan)
		select {
		case data := <-dataChan:
			userMessage := request.New(conn.GetUsername(), data)
			s.UserRequests <- userMessage
			continue
		case err := <-conn.ErrorChan:
			if err != nil {
				s.ConnectionFail(conn)
				logger.Errorf("Failed to read data from user: %v", err)
			}
			return
		}
	}
}

func (s *Server) GetUserList() []string {
	users := make([]string, 0, len(s.Connections))
	for user, _ := range s.Connections {
		users = append(users, user)
	}
	return users
}

func (s *Server) HandleRequests() {
	for {
		req := <-s.UserRequests
		logger.Infow("Got a new message", "username", req.Username, "message", req.Message, "type", req.Type)

		switch req.Type {
		case constants.FETCH_ATTENDEES_TYPE:
			go server.FetchAttendees(req.Username, s.GetUserList(), s)
		case constants.PUBLIC_MESSAGE_TYPE:
			go func() {
				err := server.SendPublicMessage(req.Username, req.Message, s)
				if err != nil {
					logger.Errorf("Error sending public message: %v", err)
				}
			}()
		case constants.PRIVATE_MESSAGE_TYPE:
			go func() {
				err := server.SendPrivateMessage(req.Username, req.Message, s)
				if err != nil {
					logger.Errorf("Error sending private message: %v", err)
				}
			}()
		case constants.BYE_MESSAGE_TYPE:
			s.Leave(req.Username)
		default:
			logger.Warnw("Unknown request", "request", req.Message)
			continue

		}
	}
}
