package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Manager struct {
	*gorm.DB
}

func NewManager(dsn string) (*Manager, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Manager{db}, nil
}

func (m *Manager) Migrate() error {
	err := m.Transaction(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(&User{})
		return err

	})
	return err
}
