package sessions

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	id        string
	csrfToken string
	data      map[string]interface{}
	expiredAt time.Time
}

func (d *Data) Id() string {
	return d.id
}

func (d *Data) CSRFToken() string {
	return d.csrfToken
}

func (d *Data) Get(key string) (interface{}, bool) {
	data, ok := d.data[key]
	return data, ok
}

// Setiap set harus disertai dengan save
func (d *Data) Set(key string, value interface{}) {
	d.data[key] = value
}

// Setiap delete harus disertai dengan save
func (d *Data) Delete(key string) {
	delete(d.data, key)
}

// Setiap clear harus disertai dengan save
func (d *Data) Clear() {
	for key := range d.data {
		delete(d.data, key)
	}
}

func (d *Data) RegenerateId() {
	d.id = uuid.NewString()
}

func (d *Data) Invalidate() {
	d.id = uuid.NewString()
	d.data = make(map[string]interface{})
}

func (d *Data) RegenerateCSRFToken() {
	d.csrfToken = uuid.NewString()
}

func (d *Data) Data() map[string]interface{} {
	return d.data
}

func (d *Data) ExpiredAt() time.Time {
	return d.expiredAt
}

func NewEmptyData(maxAge int64) *Data {

	return &Data{
		id:        uuid.NewString(),
		csrfToken: uuid.NewString(),
		data:      make(map[string]interface{}),
		expiredAt: getExpirationFromMaxAge(maxAge),
	}
}

func NewData(id string, csrfToken string, data map[string]interface{}, expiredAt time.Time) *Data {
	return &Data{
		id:        id,
		csrfToken: csrfToken,
		data:      data,
		expiredAt: expiredAt,
	}
}

func getExpirationFromMaxAge(maxAge int64) time.Time {
	return time.Now().Add(time.Minute * time.Duration(maxAge))
}
