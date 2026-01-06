package initialization

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/InstayPMS/backend/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var allModels = []any{}

type DB struct {
	Gorm *gorm.DB
	sql  *sql.DB
}

func InitPostgreSQL(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.SSLMode,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	gDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	if err := runAutoMigrations(gDB); err != nil {
		return nil, err
	}

	sqlDB, err := gDB.DB()
	if err != nil {
		return nil, err
	}

	return &DB{
		gDB,
		sqlDB,
	}, nil
}

func (d *DB) Close() {
	_ = d.sql.Close()
}

func runAutoMigrations(db *gorm.DB) error {
	return db.AutoMigrate(allModels...)
}
