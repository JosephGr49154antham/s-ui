package inbound

import (
	"errors"
	"s-ui/database"

	"gorm.io/gorm"
)

type Inbound struct {
	gorm.Model
	Tag      string `gorm:"uniqueIndex;not null" json:"tag"`
	Protocol string `gorm:"not null" json:"protocol"`
	Address  string `json:"address"`
	Port     uint16 `gorm:"not null" json:"port"`
	Enable   bool   `gorm:"default:true" json:"enable"`
}

func AutoMigrate() error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}
	return db.AutoMigrate(&Inbound{})
}

func Create(inbound *Inbound) error {
	if inbound.Tag == "" {
		return errors.New("tag is required")
	}
	if inbound.Protocol == "" {
		return errors.New("protocol is required")
	}
	if inbound.Port == 0 {
		return errors.New("port is required")
	}
	return database.GetDB().Create(inbound).Error
}

func GetAll() ([]Inbound, error) {
	var inbounds []Inbound
	err := database.GetDB().Find(&inbounds).Error
	return inbounds, err
}

func GetByTag(tag string) (*Inbound, error) {
	var inbound Inbound
	err := database.GetDB().Where("tag = ?", tag).First(&inbound).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &inbound, err
}

func Delete(tag string) error {
	result := database.GetDB().Where("tag = ?", tag).Delete(&Inbound{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("inbound not found")
	}
	return nil
}
