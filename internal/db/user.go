package db

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username      string    `json:"username" gorm:"unique"`
	Password      string    `json:"password"`
	SessionKey    string    `json:"session_key"`
	Admin         bool      `json:"admin"`
	LastLoginDate time.Time `json:"last_login_date"`
}

func (u User) UpdateLoginDate() {
	u.LastLoginDate = time.Now()
	manager.Save(u)
}

func UpdateUserLoginDate(username string) error {
	u, err := manager.GetUserByUsername(username)
	if err == nil {
		return err
	}
	u.LastLoginDate = time.Now()
	manager.Save(u)
	return nil
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
