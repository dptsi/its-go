package contracts

import "bitbucket.org/dptsi/go-framework/database"

type DatabaseService interface {
	GetDatabase(name string) *database.Database
	GetDefault() *database.Database
}
