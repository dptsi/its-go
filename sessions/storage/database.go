package storage

import (
	"context"
	"time"

	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/database"
	"bitbucket.org/dptsi/go-framework/sessions"
)

const TableName = "sessions"

type DatabaseData struct {
	Id        string                 `gorm:"primaryKey"`
	Data      map[string]interface{} `gorm:"serializer:json"`
	ExpiredAt time.Time              `gorm:"index"`
	CSRFToken string
}

func (DatabaseData) TableName() string {
	return TableName
}

type Database struct {
	db *database.Database
}

func NewDatabase(db *database.Database) *Database {
	return &Database{db}
}

func (g *Database) Get(ctx context.Context, id string) (contracts.SessionData, error) {
	var data DatabaseData
	if err := g.db.Table(TableName).First(&data, "id = ?", id).Error; err != nil {
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
	return g.db.Table(TableName).Save(&DatabaseData{data.Id(), data.Data(), data.ExpiredAt(), data.CSRFToken()}).Error
}

func (g *Database) Delete(ctx context.Context, id string) error {
	return g.db.Table(TableName).Delete(&DatabaseData{}, "id = ?", id).Error
}
