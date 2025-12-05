package database

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgresConnection creates a new PostgreSQL database connection
func NewPostgresConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Course{},
		&entity.Lesson{},
	)
}

