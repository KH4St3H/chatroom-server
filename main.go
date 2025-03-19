package main

import (
	"github.com/kh4st3h/chatroom-server/internal/config"
	"github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/log"
	"github.com/kh4st3h/chatroom-server/internal/server"
)

func main() {
	logger := log.NewLogger()
	sugar := logger.Sugar()
	sugar.Info("Loading config")
	cfg, err := config.NewConfig()
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Info("Connecting to database")
	dbManager, err := db.NewManager(cfg.DatabaseDSN)
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Info("Connected to database")

	sugar.Info("Migrating database")
	err = dbManager.Migrate()
	if err != nil {
		sugar.Fatalf("error migrating database: %v", err)
		return
	}

	srv := server.NewServer(cfg)
	err = srv.Run()
	if err != nil {
		sugar.Fatal(err)
		return
	}
}
