package user

import (
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/crypto"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/log"
)

func Register(username, password string) error {
	logger := log.NewLogger().Sugar()
	dbManager := db.GetManager()

	exists := dbManager.CheckUserExists(username)
	if exists {
		logger.Errorf("User already exists")
		return errors.New("user already exists")
	}
	encodedPassword := crypto.Base64Encode(crypto.Sha1HashData([]byte(password)))
	err := dbManager.CreateUser(username, string(encodedPassword))
	if err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return errors.New("failed to create user")
	}
	logger.Infow("user created", "username", username)
	return nil
}
