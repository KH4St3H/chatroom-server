package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/crypto"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"io"
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

func HandleRegistration(responseWriter io.Writer, ctx context.Context) {
	logger := log.NewLogger().Sugar()
	dbManager := db.GetManager()
	data := ctx.Value("data").([]byte)
	username, password, err := extractUsernamePassword(data)
	if err != nil {
		logger.Info("Failed to extract username and password from user input")
		_, _ = responseWriter.Write([]byte("Failed to extract username and password"))
		return
	}
	exists := dbManager.CheckUserExists(username)
	if exists {
		logger.Infow("user already exists", "username", username)
		_, _ = responseWriter.Write([]byte("user already exists"))
		return
	}
	encodedPassword := crypto.Base64Encode(crypto.Sha1HashData([]byte(password)))
	err = dbManager.CreateUser(username, string(encodedPassword))
	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		_, _ = responseWriter.Write([]byte("Failed to create user"))
		return
	}
	logger.Infow("user created", "username", username)
	_, _ = responseWriter.Write([]byte("Account created"))
}
