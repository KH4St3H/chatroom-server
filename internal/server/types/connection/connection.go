package connection

import (
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/crypto"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"go.uber.org/zap"
	"net"
)

var logger *zap.SugaredLogger

func init() {
	logger = log.NewLogger().Sugar()
}

type Conn struct {
	conn          net.Conn
	username      string
	cryptoManager crypto.Manager
	ErrorChan     chan error
}

func New(conn net.Conn) *Conn {
	return &Conn{conn: conn, ErrorChan: make(chan error)}
}

func (c *Conn) Write(data []byte) (int, error) {
	count, err := c.conn.Write(data)
	if err != nil {
		logger.Errorw("Failed to write data to client", "data", data)
	}
	return count, err
}

func (c *Conn) Read() ([]byte, error) {
	packet := make([]byte, 4096)
	count, err := c.conn.Read(packet)
	if err != nil {
		logger.Error("Failed to read data from client")
	}
	packet = packet[:count]
	return packet, err
}

func (c *Conn) ReadAndDecrypt(dataChan chan string) {
	packet, err := c.Read()

	if err != nil {
		c.ErrorChan <- err
		return
	}
	decrypted, err := c.DecryptMessage(packet)
	if err != nil {
		c.ErrorChan <- errors.Join(errors.New("failed to decrypt user data"), err)
		return
	}
	dataChan <- string(decrypted)
}

func (c *Conn) EncryptedWrite(msg []byte) error {
	encrypedMsg, err := c.EncryptMessage(msg)
	if err != nil {
		return errors.Join(errors.New("failed to encrypt user data"), err)
	}
	_, err = c.Write(encrypedMsg)

	if err != nil {
		return err
	}

	return nil
}

func (c *Conn) DecryptMessage(message []byte) ([]byte, error) {
	return c.cryptoManager.Decrypt(message)
}

func (c *Conn) EncryptMessage(message []byte) ([]byte, error) {
	return c.cryptoManager.Encrypt(message)
}

func (c *Conn) Authenticate(username string, sessionKey []byte) {
	c.username = username
	c.cryptoManager = *crypto.NewManager(crypto.AesCBC{}, crypto.Sha1HashData(sessionKey))
	DB := db.GetManager()
	user, err := DB.GetUserByUsername(username)
	if err != nil {
		logger.Error(err)
		return
	}
	user.SessionKey = string(crypto.Base64Encode(crypto.Sha1HashData(sessionKey)))
	DB.Save(user)
}

func (c *Conn) GoOnline() {
	DB := db.GetManager()
	event := db.NewEvent(c.GetUsername(), "connect", "")
	err := DB.SaveEvent(event)
	if err != nil {
		logger.Errorf("Failed to save event: %s", err)
	}
	err = db.UpdateUserLoginDate(c.GetUsername())
	if err != nil {
		logger.Errorf("Failed to update online status: %s", err)
	}
}

func (c *Conn) GoOffline() {
	c.ErrorChan <- nil
	DB := db.GetManager()
	event := db.NewEvent(c.GetUsername(), "disconnect", "")
	err := DB.SaveEvent(event)
	if err != nil {
		logger.Errorf("Failed to save event: %s", err)
	}
}

func (c *Conn) GetUsername() string {
	return c.username
}
