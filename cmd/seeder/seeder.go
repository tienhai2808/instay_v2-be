package main

import (
	"fmt"
	"log"

	"github.com/InstayPMS/backend/internal/domain/model"
	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"github.com/InstayPMS/backend/internal/infrastructure/initialization"
	"github.com/InstayPMS/backend/pkg/utils"
	"github.com/sony/sonyflake/v2"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := initialization.InitDatabase(cfg.PostgreSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	idGen, err := initialization.InitSnowFlake()
	if err != nil {
		log.Fatal(err)
	}

	if err := adminSeeder(db.Gorm, idGen, cfg.SuperUser.Username, cfg.SuperUser.Password); err != nil {
		log.Fatal(err)
	}
}

func adminSeeder(db *gorm.DB, idGen *sonyflake.Sonyflake, username, password string) error {
	var count int64
	if err := db.Model(&model.User{}).
		Where("role = ? AND department_id IS NULL", model.RoleAdmin).
		Count(&count).Error; err != nil {
		return fmt.Errorf("check admin failed: %w", err)
	}

	if count > 0 {
		log.Println("Admin already exists, skipping")
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}

	id, err := idGen.NextID()
	if err != nil {
		return fmt.Errorf("generate id failed: %w", err)
	}

	admin := &model.User{
		ID:        id,
		Username:  username,
		Email:     "admin@gmail.com",
		Role:      model.RoleAdmin,
		FirstName: "Main",
		LastName:  "Administrator",
		Phone:     "0123456789",
		Password:  hashedPassword,
		IsActive:  true,
	}
	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("create admin failed: %w", err)
	}

	log.Println("Admin created successfully")
	return nil
}
