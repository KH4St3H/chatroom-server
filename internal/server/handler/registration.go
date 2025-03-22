package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/actions/user"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"github.com/kh4st3h/chatroom-server/internal/server/types/connection"
	"regexp"
)

func extractUsernamePassword(data []byte) (string, string, error) {
	data = bytes.TrimSpace(data)
	r := regexp.MustCompile(`^Registration ([a-zA-Z][\w]+) (\S+)$`)
	if !r.Match(data) {
		return "", "", errors.New("invalid username or password")
	}
	matches := r.FindStringSubmatch(string(data))
	return matches[1], matches[2], nil

}

func HandleRegistration(conn connection.Conn, ctx context.Context) {
	logger := log.NewLogger().Sugar()
	data := ctx.Value("data").([]byte)
	username, password, err := extractUsernamePassword(data)
	if err != nil {
		logger.Info("Failed to extract username and password from user input")
		_, _ = conn.Write([]byte("Failed to extract username and password"))
		return
	}
	err = user.Register(username, password)
	if err != nil {
		conn.Write([]byte(err.Error()))
	}
	_, _ = conn.Write([]byte("Account created"))
}
