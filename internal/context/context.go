package context

import (
	"github.com/kh4st3h/chatroom-server/internal/config"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"go.uber.org/zap"
)

type Context struct {
	Logger    *zap.SugaredLogger
	Config    *config.Config
	DBManager *db.Manager
}
