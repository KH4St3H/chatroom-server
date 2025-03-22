package user

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/kh4st3h/chatroom-server/internal/crypto"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	SessionVerifyToken string `json:"session_token"`
	SessionKey         []byte `json:"session_key"`
}

type VerifyLoginRequest struct {
	Data       []byte
	Username   string
	SessionKey []byte
}

type VerifyLoginResponse struct {
	Ok       bool   `json:"ok"`
	Username string `json:"username"`
}

func Login(req LoginRequest) (*LoginResponse, error) {
	logger := log.NewLogger().Sugar()
	dbManager := db.GetManager()

	username := req.Username

	user, err := dbManager.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Info("User not found")
			return nil, errors.New("User not found")
		}
	}
	randomKey, err := crypto.GenerateRandomAESKey()
	if err != nil {
		logger.Error("Failed to generate random key")
		return nil, errors.New("internal error")
	}
	logger.Debugw("Generated random key", "key", randomKey)
	password, err := crypto.Base64Decode([]byte(user.Password))
	if err != nil {
		logger.Error("Failed to decode password")
		return nil, errors.New("internal error")
	}
	cryptoManager := crypto.NewManager(&crypto.AesCBC{}, password)

	encryptedKey, err := cryptoManager.Encrypt(randomKey)
	if err != nil {
		logger.Error("Failed to encrypt generated password")
		return nil, errors.New("internal error")
	}

	response := hex.EncodeToString(encryptedKey)

	return &LoginResponse{SessionVerifyToken: response, SessionKey: randomKey}, nil
}

func VerifyLogin(request VerifyLoginRequest) (VerifyLoginResponse, error) {
	cryptoManager := crypto.NewManager(&crypto.AesCBC{}, crypto.Sha1HashData(request.SessionKey))
	plaintext, err := cryptoManager.Decrypt(request.Data)
	if err != nil {
		return VerifyLoginResponse{Ok: false}, err
	}
	if string(plaintext) == fmt.Sprintf("Hello %s", request.Username) {
		return VerifyLoginResponse{Ok: true, Username: request.Username}, nil
	}
	return VerifyLoginResponse{Ok: false, Username: request.Username}, nil
}
