package database

import "gorm.io/gorm"

type Database = gorm.DB

var ErrRecordNotFound = gorm.ErrRecordNotFound
