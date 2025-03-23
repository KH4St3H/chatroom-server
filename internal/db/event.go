package db

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	UserRef string
	User    User `gorm:"references:username;foreignKey:UserRef;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Type    string
	Message string
}

func NewEvent(username string, eventType string, message string) Event {
	return Event{
		UserRef: username,
		Type:    eventType,
		Message: message,
	}
}

func (m *Manager) SaveEvent(event Event) error {
	return m.DB.Save(&event).Error
}
