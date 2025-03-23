package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Manager struct {
	*gorm.DB
}

var manager *Manager

func GetManager() *Manager {
	return manager
}

func NewManager(dsn string) (*Manager, error) {
	if manager != nil {
		return manager, nil
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	manager = &Manager{db}
	return manager, nil
}

func (m *Manager) Migrate() error {
	err := m.Transaction(func(tx *gorm.DB) error {
		err := tx.AutoMigrate(&User{})
		if err != nil {
			return err
		}
		err = tx.AutoMigrate(&Event{})
		return err
	})
	return err
}
