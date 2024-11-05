package storage

import (
	"context"
	"time"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/database"
	"github.com/dptsi/its-go/sessions"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
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
		db.Table(table).AutoMigrate(&DatabaseData{})
	}
	return &Database{db, table}
}

func (g *Database) Get(ctx context.Context, id string) (contracts.SessionData, error) {
	var data DatabaseData
	if err := uuid.Validate(id); err != nil {
		return nil, nil
	}
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
	err := g.db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"data":       data.Data(),
					"expired_at": data.ExpiredAt(),
					"csrf_token": data.CSRFToken(),
				},
			),
		},
	).
		Table(g.table).
		Save(&DatabaseData{
			data.Id(),
			data.Data(),
			data.ExpiredAt(),
			data.CSRFToken(),
		}).Error

	return err
}

func (g *Database) Delete(ctx context.Context, id string) error {
	return g.db.Table(g.table).Where("id = ?", id).Delete(&DatabaseData{}).Error
}
