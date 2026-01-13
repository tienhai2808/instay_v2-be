package seeder

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/sony/sonyflake/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Seeder struct {
	cfg      config.SuperUserConfig
	log      *zap.Logger
	db       *gorm.DB
	idGen    *sonyflake.Sonyflake
	userRepo repository.UserRepository
}

func NewSeeder(
	cfg config.SuperUserConfig,
	log *zap.Logger,
	db *gorm.DB,
	idGen *sonyflake.Sonyflake,
	userRepo repository.UserRepository,
) *Seeder {
	return &Seeder{
		cfg,
		log,
		db,
		idGen,
		userRepo,
	}
}

func (s *Seeder) Start() error {
	return s.adminSeeder(s.cfg.Username, s.cfg.Password)
}
