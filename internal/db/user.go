package db

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password"`
	SessionKey string `json:"session_key"`
	Admin      bool   `json:"admin"`
}

func (m *Manager) CheckUserExists(username string) bool {
	var user User
	result := m.Take(&user, "username = ?", username)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func (m *Manager) CreateUser(username string, password string) error {
	user := User{
		Username: username,
		Admin:    false,
		Password: password,
	}
	result := m.Create(&user)
	return result.Error
}

func (m *Manager) GetUserByUsername(username string) (*User, error) {
	var user User
	result := m.Take(&user, "username = ?", username)
	return &user, result.Error
}

func (m *Manager) SaveUser(user *User) error {
	return m.Save(user).Error

}
