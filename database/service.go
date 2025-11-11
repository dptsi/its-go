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

const (
	DefaultTrustServerCertificate string = "true"
	DefaultTransportEncrypt       string = "enabled"
)

type ConnectionConfig struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	Database string

	// not used on `sqlite` driver
	TrustServerCertificate string

	// not used on `sqlite` driver
	TransportEncrypt string
}

type Config struct {
	Connections map[string]ConnectionConfig
}

func NewService(cfg Config) (*Service, error) {
	/**
	https://github.com/dptsi/its-go/blob/2b8efc3f44d4ecc95a63ccd5b9d7a6a3cd0b3a30/database/service.go#L25

	Add this configuration:

		* `TrustCertificate`: option to disable forced certificate trust. default is `true`.

		* `Encrypt`: option for transport layer encryption (i assume?). default is `true.

	https://pkg.go.dev/github.com/microsoft/go-mssqldb@v1.6.0#section-readme
	*/
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

	if cfg.TransportEncrypt == "" {
		cfg.TransportEncrypt = DefaultTransportEncrypt
	}

	if cfg.TrustServerCertificate == "" {
		cfg.TrustServerCertificate = DefaultTrustServerCertificate
	}

	if cfg.TrustServerCertificate != "true" && cfg.TrustServerCertificate != "false" {
		return nil, fmt.Errorf("invalid database TrustServerCertificate configuration: %s", cfg.TrustServerCertificate)
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
		/**
		https://pkg.go.dev/github.com/microsoft/go-mssqldb@v1.6.0#section-readme

		encrypt
			strict - Data sent between client and server is encrypted E2E using TDS8.
			disable - Data send between client and server is not encrypted.
			false/optional/no/0/f - Data sent between client and server is not encrypted beyond the login packet.
			true/mandatory/yes/1/t - Data sent between client and server is encrypted.

		TrustServerCertificate
			false - Server certificate is checked. Default is false if encrypt is specified.
			true - Server certificate is not checked. Default is true if encrypt is not specified. If trust server
				   certificate is true, driver accepts any certificate presented by the server and any host name in that
				   certificate. In this mode, TLS is susceptible to man-in-the-middle attacks. This should be used only for testing.
		*/
		/**
		(comment id: its-go/database/service.go-1)
		in the making of this encrypt configuration that works between mssql and
		postgresql, we decided that on sqlserver the usage definiton of "disabled"
		and "false" is confusing. as above, "disable" is defined as "Data send
		between client and server is not encrypted", while "false" is "Data sent
		between client and server is not encrypted [[beyond the login packet]].".
		the consequences of using "false" instead of "disable" is you may get
		tls cert verify error like below:

		...
		[error] failed to initialize database, got error TLS Handshake
		failed: tls: failed to verify certificate: x509: cannot validate
		certificate for (db ip/host) because it doesn't contain any IP SANs
		...

		(comment id: its-go/database/service.go-1-2)
		"[[beyond the login packet]]" implies there's transport encryption that
		is enabled at login. if you specify TransportEncrypt as either "false"
		or "login-only", and you get tls handshake verify certificate error,
		make sure to configure TrustServerCertificate as "true".
		*/
		transportEncrypt := ""
		switch cfg.TransportEncrypt {
		case "enable", "enabled":
			transportEncrypt = "&encrypt=true"
		case "strict", "true", "mandatory", "yes", "1", "t", "optional", "disable", "false", "no", "0", "f":
			transportEncrypt = fmt.Sprintf("&encrypt=%s", cfg.TransportEncrypt)
		case "disabled":
			// refer to comment id its-go/database/service.go-1
			transportEncrypt = "&encrypt=disable"
		case "login-only":
			// refer to comment id its-go/database/service.go-1-2
			transportEncrypt = "&encrypt=false"
		default:
			return nil, fmt.Errorf("invalid database TransportEncrypt configuration: %s", cfg.TransportEncrypt)
		}

		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%s?database=%s%s&trustservercertificate=%s",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
			transportEncrypt,
			cfg.TrustServerCertificate,
		)
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQL Server connection error: %w", err)
		}
		return db, nil

	case "postgres":
		/**
		set sslmode accordingly
			https://www.postgresql.org/docs/current/libpq-ssl.html
			https://www.postgresql.org/docs/current/libpq-ssl.html#LIBPQ-SSL-SSLMODE-STATEMENTS

		sslmode		Eavesdropping 	MITM		Statement
					protection		protection
		disable 	No 				No 			I don't care about security, and I don't want to pay the overhead of encryption.
		allow 		Maybe 			No 			I don't care about security, but I will pay the overhead of encryption if the server insists on it.
		prefer 		Maybe 			No 			I don't care about encryption, but I wish to pay the overhead of encryption if the server supports it.
		require 	Yes 			No 			I want my data to be encrypted, and I accept the overhead. I trust that the network will make sure I
												always connect to the server I want.
		verify-ca 	Yes 			Depends on	I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server that I trust.
									CA policy
		verify-full Yes 			Yes			I want my data encrypted, and I accept the overhead. I want to be sure that I connect to a server I trust,
												and that it's the one I specify.
		*/
		sslmode := ""
		switch cfg.TransportEncrypt {
		case "strict":
			sslmode = "verify-full"
		case "enable", "enabled", "true", "mandatory", "yes", "1", "t":
			if cfg.TrustServerCertificate == "true" {
				sslmode = "verify-ca"
				break
			}
			sslmode = "require"
		case "optional":
			if cfg.TrustServerCertificate == "true" {
				sslmode = "verify-ca"
				break
			}
			sslmode = "prefer"
		case "disable", "disabled", "false", "no", "0", "f":
			sslmode = "disable"
		default:
			return nil, fmt.Errorf("invalid database TrustServerCertificate configuration: %s", cfg.TransportEncrypt)
		}

		params := []string{
			fmt.Sprintf("host=%s", cfg.Host),
			fmt.Sprintf("user=%s", cfg.User),
			fmt.Sprintf("password=%s", cfg.Password),
			fmt.Sprintf("dbname=%s", cfg.Database),
			fmt.Sprintf("sslmode=%s", sslmode),
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
