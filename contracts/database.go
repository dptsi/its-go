package contracts

import "github.com/dptsi/its-go/database"

type DatabaseService interface {
	GetDatabase(name string) *database.Database
	GetDefault() *database.Database
}
