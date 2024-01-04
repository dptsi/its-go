package storage

import (
	"context"
	"log"
	"time"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/database"
	"github.com/dptsi/its-go/sessions"
	"gorm.io/gorm"
)

type DatabaseData struct {
	Id        string                 `gorm:"primaryKey"`
	Data      map[string]interface{} `gorm:"serializer:json"`
	ExpiredAt time.Time              `gorm:"index"`
	CSRFToken string
}

type Database struct {
	db    *database.Database
	table string
}

func NewDatabase(db *database.Database, table string, autoMigrate bool) *Database {
	if autoMigrate {
		log.Printf("Auto migrate sessions table with name %s", table)
		db.Table(table).AutoMigrate(&DatabaseData{})
		log.Printf("Table %s successfully migrated", table)
	}
	return &Database{db, table}
}

func (g *Database) Get(ctx context.Context, id string) (contracts.SessionData, error) {
	var data DatabaseData
	if err := g.db.Table(g.table).Where("id = ?", id).First(&data).Error; err != nil {
		if err == database.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if data.ExpiredAt.Before(time.Now()) {
		return nil, nil
	}

	sess := sessions.NewData(id, data.CSRFToken, data.Data, data.ExpiredAt)
	return sess, nil
}

func (g *Database) Save(ctx context.Context, data contracts.SessionData) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(g.table).Delete(&DatabaseData{}, "id = ?", data.Id()).Error; err != nil {
			return err
		}
		if err := tx.Table(g.table).Save(&DatabaseData{data.Id(), data.Data(), data.ExpiredAt(), data.CSRFToken()}).Error; err != nil {
			return err
		}

		return nil
	})
}
