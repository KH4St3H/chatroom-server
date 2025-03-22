package connection

import (
	"errors"
	"github.com/kh4st3h/chatroom-server/internal/crypto"
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
}

func New(conn net.Conn) *Conn {
	return &Conn{conn: conn}
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

func (c *Conn) ReadAndDecrypt() (string, error) {
	packet, err := c.Read()

	if err != nil {
		return "", err
	}
	decrypted, err := c.DecryptMessage(packet)
	if err != nil {
		return "", errors.Join(errors.New("Failed to decrypt user data"), err)
	}
	return string(decrypted), nil
}

func (c *Conn) DecryptMessage(message []byte) ([]byte, error) {
	return c.cryptoManager.Decrypt(message)
}

func (c *Conn) Authenticate(username string, sessionKey []byte) {
	c.username = username
	c.cryptoManager = *crypto.NewManager(crypto.AesCBC{}, crypto.Sha1HashData(sessionKey))
}

func (c *Conn) GoOffline() {

}

func (c *Conn) GetUsername() string {
	return c.username
}
