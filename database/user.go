package database

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// User represents an application user stored in the database.
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"uniqueIndex;not null"     json:"username"`
	Password  string         `gorm:"not null"                 json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                    json:"-"`
}

// CreateUser inserts a new user into the database.
func CreateUser(username, hashedPassword string) (*User, error) {
	if username == "" {
		return nil, errors.New("username must not be empty")
	}
	u := &User{Username: username, Password: hashedPassword}
	result := db.Create(u)
	return u, result.Error
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*User, error) {
	var u User
	result := db.Where("username = ?", username).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &u, result.Error
}

// DeleteUser soft-deletes a user by ID.
func DeleteUser(id uint) error {
	return db.Delete(&User{}, id).Error
}
