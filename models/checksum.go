package models

import "gorm.io/gorm"

type ChecksumData struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	ChecksumValue string `gorm:"column:checksum_value;unique"`
}
