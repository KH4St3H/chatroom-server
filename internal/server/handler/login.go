package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/kh4st3h/chatroom-server/internal/actions/user"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"github.com/kh4st3h/chatroom-server/internal/server/types/connection"
	"regexp"
)

func extractUsername(data []byte) (string, error) {
	data = bytes.TrimSpace(data)
	r := regexp.MustCompile(`^Login ([a-zA-Z]\w+)$`)
	if !r.Match(data) {
		return "", errors.New("invalid username or password")
	}
	matches := r.FindStringSubmatch(string(data))
	return matches[1], nil

}

func HandleLogin(conn *connection.Conn, ctx context.Context) bool {
	logger := log.NewLogger().Sugar()
	data := ctx.Value("data").([]byte)

	username, err := extractUsername(data)
	if err != nil {
		logger.Info("Failed to extract username from request")
		_, _ = conn.Write([]byte("Failed to extract username from request"))
		return false
	}
	loginRequest := user.LoginRequest{Username: username}
	loginResponse, err := user.Login(loginRequest)
	if err != nil {
		_, err := conn.Write([]byte(err.Error()))
		if err != nil {
			return false
		}
		return false
	}
	_, err = conn.Write([]byte(fmt.Sprintf("Key %s", loginResponse.SessionVerifyToken)))
	if err != nil {
		return false
	}
	newCtx := context.WithValue(ctx, "username", username)
	newCtx = context.WithValue(newCtx, "sessionKey", loginResponse.SessionKey)

	success := VerifyLogin(*conn, newCtx)
	if success {
		conn.Authenticate(username, loginResponse.SessionKey)
		conn.Write([]byte(fmt.Sprintf("Hi %s, welcome to the chatroom.", conn.GetUsername())))
		return true
	}
	return false
}

func VerifyLogin(conn connection.Conn, ctx context.Context) bool {
	logger := log.NewLogger().Sugar()
	data, err := conn.Read()
	if err != nil {
		logger.Errorf("Failed to read login verification response: %v", err)
		return false
	}
	response, err := user.VerifyLogin(user.VerifyLoginRequest{
		Data: data, SessionKey: ctx.Value("sessionKey").([]byte), Username: ctx.Value("username").(string),
	})
	if err != nil {
		logger.Errorf("Failed to verify login: %v", err)
		conn.Write([]byte("Failed to authorize"))
		return false
	}

	if response.Ok != true {
		logger.Error("Failed to verify login")
		conn.Write([]byte("Failed to authorize"))
		return false
	}
	return true
}
