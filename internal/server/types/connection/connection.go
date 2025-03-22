package connection

import (
	"github.com/kh4st3h/chatroom-server/internal/log"
	"go.uber.org/zap"
	"net"
)

var logger *zap.SugaredLogger

func init() {
	logger = log.NewLogger().Sugar()
}

type Conn struct {
	conn net.Conn
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
