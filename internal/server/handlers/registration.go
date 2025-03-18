package handlers

import (
	"bytes"
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/context"
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

func HandleRegistration(ctx *context.Context, data []byte) error {
	username, password, err := extractUsernamePassword(data)
	if err != nil {
		return err
	}
	exists := ctx.DBManager.CheckUserExists(username)
	if exists {
		return errors.New("user already exists")
	}
	err = ctx.DBManager.CreateUser(username, password)
	if err != nil {
		ctx.Logger.Errorf("Failed to create user: %v", err)
		return errors.New("could not create user")
	}
	return nil
}
