package database

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const DefaultDatabaseName = "default"

type Service struct {
	databases map[string]*Database
}

type ConnectionConfig struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Database string
	Timezone string
}

type Config struct {
	Connections map[string]ConnectionConfig
}

func NewService(cfg Config) (*Service, error) {
	databases := make(map[string]*Database, len(cfg.Connections))
	for name, cfg := range cfg.Connections {
		if name == "" {
			name = DefaultDatabaseName
		}
		db, err := createConnection(cfg)
		if err != nil {
			return nil, fmt.Errorf("database service: new service: error creating database with name \"%s\": %w", name, err)
		}
		databases[name] = db
	}

	return &Service{
		databases: databases,
	}, nil
}

func createConnection(cfg ConnectionConfig) (*Database, error) {
	if cfg.Driver == "" {
		return nil, fmt.Errorf("database driver is empty, supported drivers are [sqlite, sqlserver, postgres]")
	}

	// set default timezone if not provided
	if cfg.Timezone == "" {
		cfg.Timezone = "Asia/Jakarta"
	}

	switch cfg.Driver {
	case "sqlite":
		// Contoh penggunaan adapter GORM dengan SQLite
		db, err := gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQLite connection error: %w", err)
		}
		return db, nil
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQL Server connection error: %w", err)
		}
		return db, nil
	case "postgres":
		params := []string{
			fmt.Sprintf("host=%s", cfg.Host),
			fmt.Sprintf("user=%s", cfg.User),
			fmt.Sprintf("password=%s", cfg.Password),
			fmt.Sprintf("dbname=%s", cfg.Database),
			fmt.Sprintf("TimeZone=%s", cfg.Timezone),
			"sslmode=disable",
		}

		if cfg.Port != "" {
			params = append(params, fmt.Sprintf("port=%s", cfg.Port))
		}

		dsn := strings.Join(params, " ")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("PostgreSQL connection error: %w", err)
		}
		return db, nil
	default:
		return nil, fmt.Errorf("unknown database driver %s, supported drivers are [sqlite, sqlserver, postgres]", cfg.Driver)
	}

}

func (m *Service) GetDatabase(name string) *Database {
	return m.databases[name]
}

func (m *Service) GetDefault() *Database {
	return m.databases[DefaultDatabaseName]
}
