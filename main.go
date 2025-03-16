package main

import (
	"github.com/kh4st3h/chatroom-server/internal/config"
	db2 "github.com/kh4st3h/chatroom-server/internal/db"
	"github.com/kh4st3h/chatroom-server/internal/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Info("Loading config")
	cfg, err := config.NewConfig()
	if err != nil {
		sugar.Fatal(err)
	}
	sugar.Info("Connecting to database")
	dbManager, err := db2.NewManager(cfg.DatabaseDSN)
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

	srv := server.NewServer(cfg, sugar, dbManager)
	err = srv.Run()
	if err != nil {
		sugar.Fatal(err)
		return
	}
}
