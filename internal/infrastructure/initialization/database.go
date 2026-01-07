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

type Database struct {
	Gorm *gorm.DB
	sql  *sql.DB
}

func InitDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.PostgreSQL.Host,
		cfg.PostgreSQL.DBName,
		cfg.PostgreSQL.User,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.SSLMode,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "[DB] ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	if err := runAutoMigrations(pg); err != nil {
		return nil, err
	}

	sql, err := pg.DB()
	if err != nil {
		return nil, err
	}

	return &Database{
		pg,
		sql,
	}, nil
}

func (d *Database) Close() {
	_ = d.sql.Close()
}

var allModels = []any{}

func runAutoMigrations(db *gorm.DB) error {
	return db.AutoMigrate(allModels...)
}
