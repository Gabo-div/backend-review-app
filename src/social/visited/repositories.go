package visited

import (
	"errors"

	"github.com/NetKBs/backend-reviewapp/config"
	"github.com/NetKBs/backend-reviewapp/src/schema"
	"gorm.io/gorm"
)

func GetVisitedPlacesByUserIdRepository(userId uint, limit int, cursor uint) ([]uint, error) {
	db := config.DB
	var user schema.User
	var visitedPlaces []schema.Place

	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	query := db.Model(&user).Preload("VisitedPlaces", func(db *gorm.DB) *gorm.DB {
		if cursor != 0 {
			return db.Order("id DESC").Where("id < ?", cursor).Limit(limit)
		}
		return db.Order("id DESC").Limit(limit)
	})

	if err := query.Find(&user).Error; err != nil {
		return nil, err
	}

	visitedPlaces = user.VisitedPlaces

	var visitedPlaceIDs []uint
	for _, place := range visitedPlaces {
		visitedPlaceIDs = append(visitedPlaceIDs, place.ID)
	}

	return visitedPlaceIDs, nil
}

func GetVisitedCountRepository(userId uint) (visitedCount uint, err error) {
	db := config.DB
	var user schema.User

	if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
		return visitedCount, err
	}

	visitedCount = uint(db.Model(&user).Association("VisitedPlaces").Count())
	return visitedCount, nil
}

func GetVisitorsCountRepository(placeID uint) (visitorsCount uint, err error) {
	db := config.DB
	var count int64
	err = db.Table("place_visitors").Where("place_id = ?", placeID).Count(&count).Error
	return uint(count), err //visitorsCount, err
}

func CreateVisitedPlaceRepository(userID, placeID uint) error {
	db := config.DB

	var user schema.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	var place schema.Place
	if err := db.Where("id = ?", placeID).First(&place).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("place not found")
		}
		return err
	}

	if err := db.Model(&user).Association("VisitedPlaces").Append(&place); err != nil {
		return err
	}

	return nil
}

func DeleteVisitedPlaceRepository(userID, placeID uint) error {
	db := config.DB

	var user schema.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	var place schema.Place
	if err := db.Where("id = ?", placeID).First(&place).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("place not found")
		}
		return err
	}

	if err := db.Model(&user).Association("VisitedPlaces").Delete(&place); err != nil {
		return err
	}

	return nil
}
