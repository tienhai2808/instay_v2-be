package seeder

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"go.uber.org/zap"
)

func (s *Seeder) adminSeeder(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := s.userRepo.ExistsActiveAdmin(ctx)
	if err != nil {
		s.log.Error("check active admin failed", zap.Error(err))
		return err
	}
	if exists {
		s.log.Info("Admin already exists, skipping")
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		s.log.Error("hash password failed", zap.Error(err))
		return err
	}

	id, err := s.idGen.NextID()
	if err != nil {
		s.log.Error("generate admin id failed", zap.Error(err))
		return err
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
	if err := s.userRepo.Create(ctx, admin); err != nil {
		s.log.Error("create admin failed", zap.Error(err))
		return err
	}

	s.log.Info("Admin created successfully")
	return nil
}
