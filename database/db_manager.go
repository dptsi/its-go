package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const DefaultDatabaseName = "default"

type Manager struct {
	databases map[string]*Database
}

type Config struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func NewManager() *Manager {
	return &Manager{
		databases: make(map[string]*Database),
	}
}

func (m *Manager) AddDatabase(name string, cfg Config) error {
	log.Printf("Adding database %s...\n", name)
	if _, exists := m.databases[name]; exists {
		return fmt.Errorf("database %s already exists", name)
	}
	if name == "" {
		name = DefaultDatabaseName
	}
	switch cfg.Driver {
	case "sqlite":
		// Contoh penggunaan adapter GORM dengan SQLite
		log.Println("Connecting to SQLite database...")
		db, err := gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("SQLite connection error: %w", err)
		}
		log.Println("Successfully connected to SQLite database!")
		m.databases[name] = NewDatabase(db)
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		log.Println("Connecting to SQL Server database...")
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("SQL Server connection error: %w", err)
		}
		log.Println("Successfully connected to SQL Server database!")
		m.databases[name] = NewDatabase(db)
	case "postgres":
		dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		log.Println("Connecting to PostgreSQL database...")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("PostgreSQL connection error: %w", err)
		}
		log.Println("Successfully connected to PostgreSQL database!")
		m.databases[name] = NewDatabase(db)
	default:
		return fmt.Errorf("unknown database driver %s", cfg.Driver)
	}

	return nil
}

func (m *Manager) GetDatabase(name string) *Database {
	return m.databases[name]
}

func (m *Manager) GetDefault() *Database {
	return m.databases[DefaultDatabaseName]
}
