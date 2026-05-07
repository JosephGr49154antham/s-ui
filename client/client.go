package client

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Client struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Enable   bool   `gorm:"default:true" json:"enable"`
	InboundTag string `gorm:"not null" json:"inbound_tag"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Client{})
}

func Create(db *gorm.DB, c *Client) error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.Email == "" {
		return errors.New("email is required")
	}
	if c.InboundTag == "" {
		return errors.New("inbound_tag is required")
	}
	// Normalize email to lowercase before saving
	c.Email = strings.ToLower(c.Email)
	return db.Create(c).Error
}

func GetAll(db *gorm.DB) ([]Client, error) {
	var clients []Client
	err := db.Find(&clients).Error
	return clients, err
}

func GetByEmail(db *gorm.DB, email string) (*Client, error) {
	var c Client
	// Normalize email to lowercase for consistent lookups
	err := db.Where("email = ?", strings.ToLower(email)).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func Delete(db *gorm.DB, id uint) error {
	result := db.Delete(&Client{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("client not found")
	}
	return nil
}
