package adapters

import (
	"context"
	"time"

	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"gorm.io/gorm"
)

const TableName = "sessions"

type GormData struct {
	Id        string                 `gorm:"primaryKey"`
	Data      map[string]interface{} `gorm:"serializer:json"`
	ExpiredAt time.Time              `gorm:"index"`
	CSRFToken string
}

func (GormData) TableName() string {
	return TableName
}

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{db}
}

func (g *Gorm) Get(ctx context.Context, id string) (*sessions.Data, error) {
	var data GormData
	if err := g.db.Table(TableName).First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
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

func (g *Gorm) Save(ctx context.Context, data *sessions.Data) error {
	return g.db.Table(TableName).Save(&GormData{data.Id(), data.Data(), data.ExpiredAt(), data.CSRFToken()}).Error
}

func (g *Gorm) Delete(ctx context.Context, id string) error {
	return g.db.Table(TableName).Delete(&GormData{}, "id = ?", id).Error
}
