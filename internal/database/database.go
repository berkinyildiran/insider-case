package database

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Database struct {
	config *Config

	context context.Context
	DB      *gorm.DB
	WG      sync.WaitGroup
}

func NewDatabase(config *Config, context context.Context) *Database {
	return &Database{
		config: config,

		context: context,
	}
}

func (d *Database) Close() error {
	d.WG.Wait()

	db, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return db.Close()
}

func (d *Database) Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", d.config.Host, d.config.User, d.config.Password, d.config.Name, d.config.Port)
	dialector := postgres.Open(dsn)

	db, err := gorm.Open(dialector)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	d.DB = db

	return nil
}

func (d *Database) Migrate(models ...any) error {
	err := d.DB.WithContext(d.context).AutoMigrate(models...)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
