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

type Service struct {
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

func NewService(databasesConfig map[string]Config) (*Service, error) {
	databases := make(map[string]*Database, len(databasesConfig))
	for name, cfg := range databasesConfig {
		log.Printf("Adding database %s...\n", name)
		if name == "" {
			name = DefaultDatabaseName
		}
		db, err := createDatabase(cfg)
		if err != nil {
			return nil, fmt.Errorf("database service: new service: error creating database with name \"%s\": %w", name, err)
		}
		databases[name] = db
	}

	return &Service{
		databases: databases,
	}, nil
}

func createDatabase(cfg Config) (*Database, error) {
	switch cfg.Driver {
	case "sqlite":
		// Contoh penggunaan adapter GORM dengan SQLite
		log.Println("Connecting to SQLite database...")
		db, err := gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQLite connection error: %w", err)
		}
		log.Println("Successfully connected to SQLite database!")
		return db, nil
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		log.Println("Connecting to SQL Server database...")
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQL Server connection error: %w", err)
		}
		log.Println("Successfully connected to SQL Server database!")
		return db, nil
	case "postgres":
		dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		log.Println("Connecting to PostgreSQL database...")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("PostgreSQL connection error: %w", err)
		}
		log.Println("Successfully connected to PostgreSQL database!")
		return db, nil
	default:
		return nil, fmt.Errorf("unknown database driver %s", cfg.Driver)
	}

}

func (m *Service) GetDatabase(name string) *Database {
	return m.databases[name]
}

func (m *Service) GetDefault() *Database {
	return m.databases[DefaultDatabaseName]
}
