package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
)

type LoadTrackingRepository struct {
	db *gorm.DB
}

func NewLoadTrackingRepository() *LoadTrackingRepository {
	return &LoadTrackingRepository{
		db: models.DB,
	}
}
