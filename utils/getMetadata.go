package utils

import (
	"errors"

	"github.com/relumini/shortdl/models"
	"gorm.io/gorm"
)

func GetMetadata(db *gorm.DB, checksum string) (models.ChecksumData, error) {
	var checkSum models.ChecksumData
	result := db.Where("checksum_value = ?", checksum).First(&checkSum)
	if result.RowsAffected == 0 {
		return models.ChecksumData{}, errors.New("not found")
	}
	return checkSum, nil
}
