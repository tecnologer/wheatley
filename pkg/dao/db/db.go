package db

import (
	"fmt"

	"github.com/tecnologer/wheatley/pkg/models"
	"github.com/tecnologer/wheatley/pkg/utils/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connection struct {
	*gorm.DB
	hasTransaction bool
}

func NewConnection(config *Config) (*Connection, error) {
	log.Infof("connecting to DB %s", config.DBName)

	fileDBName, err := config.fileDBName()
	if err != nil {
		return nil, fmt.Errorf("create DB file name: %w", err)
	}

	gormDB, err := gorm.Open(sqlite.Open(fileDBName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	log.Infof("connection established to DB %s", config.DBName)

	err = Migrate(gormDB)
	if err != nil {
		return nil, fmt.Errorf("migrate models: %w", err)
	}

	log.Infof("migrations applied to DB %s", config.DBName)

	return &Connection{
		DB: gormDB,
	}, nil
}

func (c *Connection) BeginTransaction() error {
	if c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Begin()
	c.hasTransaction = true

	log.Debug("transaction started")

	return nil
}

func (c *Connection) Commit() error {
	if !c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Commit()
	c.hasTransaction = false

	log.Debug("transaction committed")

	return nil
}

func (c *Connection) Rollback() error {
	if !c.hasTransaction {
		return nil
	}

	c.DB = c.DB.Rollback()
	c.hasTransaction = false

	log.Debug("transaction rolled back")

	return nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Notification{},
	)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
