package seeders

import (
	"log"

	"github.com/ardhisparahita/ecommerce-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	var count int64

	db.Model(&domain.User{}).Where("email = ?", "admin@gmail.com").Count(&count)
	if count > 0 {
		log.Println("Admin already exist")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	admin := domain.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: string(hashedPassword),
		Role:     "ADMIN",
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}
