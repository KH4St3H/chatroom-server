package handler

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/crypto"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"gorm.io/gorm"
	"io"
	"regexp"
)

func extractUsername(data []byte) (string, error) {
	data = bytes.TrimSpace(data)
	r := regexp.MustCompile(`^Login ([a-zA-Z][\w]+)$`)
	if !r.Match(data) {
		return "", errors.New("invalid username or password")
	}
	matches := r.FindStringSubmatch(string(data))
	return matches[1], nil

}

func HandleLogin(responseWriter io.Writer, ctx context.Context) bool {
	logger := log.NewLogger().Sugar()
	dbManager := db.GetManager()
	data := ctx.Value("data").([]byte)

	username, err := extractUsername(data)
	if err != nil {
		logger.Info("Failed to extract username from message")
		_, _ = responseWriter.Write([]byte("Failed to extract username from message"))
		return false
	}
	user, err := dbManager.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Info("User not found")
			_, _ = responseWriter.Write([]byte("User not found"))
		}
		return false
	}
	randomKey, err := crypto.GenerateRandomAESKey()
	if err != nil {
		logger.Error("Failed to generate random key")
		_, _ = responseWriter.Write([]byte("internal error"))
		return false
	}
	cryptoManager := crypto.NewManager(&crypto.AesCBC{})
	password, err := crypto.Base64Decode([]byte(user.Password))
	if err != nil {
		logger.Error("Failed to decode password")
		_, _ = responseWriter.Write([]byte("internal error"))
	}

	encryptedKey, err := cryptoManager.Encrypt(randomKey, password)
	if err != nil {
		logger.Error("Failed to encrypt generated password")
		return false
	}
	user.SessionKey = hex.EncodeToString(randomKey)
	err = dbManager.SaveUser(user)
	if err != nil {
		logger.Errorf("Failed to save user session key: %v", err)
		_, err = responseWriter.Write([]byte("internal error"))
		return false
	}
	response := hex.EncodeToString(encryptedKey)

	_, err = responseWriter.Write([]byte("Login " + response))
	return err == nil
}
