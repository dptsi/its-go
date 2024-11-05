package providers

import (
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/dptsi/its-go/activitylog"
	"github.com/dptsi/its-go/app"
	"github.com/dptsi/its-go/auth"
	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/database"
	"github.com/dptsi/its-go/event"
	_firestore "github.com/dptsi/its-go/firestore"
	"github.com/dptsi/its-go/http/middleware"
	"github.com/dptsi/its-go/logging"
	"github.com/dptsi/its-go/module"
	"github.com/dptsi/its-go/sessions"
	"github.com/dptsi/its-go/sessions/storage"
	"github.com/dptsi/its-go/web"
)

func LoadProviders(application contracts.Application) error {
	config := application.Config()
	// cryptConfig, ok := config["crypt"].(crypt.Config)
	// if !ok {
	// 	return fmt.Errorf("crypt config is not available")
	// }
	dbConfig, ok := config["database"].(database.Config)
	if !ok {
		return fmt.Errorf("database config is not available")
	}
	middlewareConfig, ok := config["middleware"].(middleware.Config)
	if !ok {
		return fmt.Errorf("middleware config is not available")
	}
	sessionsConfig, ok := config["sessions"].(sessions.Config)
	if !ok {
		return fmt.Errorf("sessions config is not available")
	}
	webConfig, ok := config["web"].(web.Config)
	if !ok {
		return fmt.Errorf("web config is not available")
	}

	app.Bind(application, "firestore.client", func(application contracts.Application) (*firestore.Client, error) {
		config, ok := config["firestore"].(_firestore.Config)
		if !ok {
			return nil, fmt.Errorf("firestore config is not available")
		}
		return firestore.NewClient(application.Context(), config.ProjectId)
	})

	app.Bind(application, "activity_log.service", func(application contracts.Application) (contracts.ActivityLogService, error) {
		logger := app.MustMake[contracts.LoggingService](application, "logging.service")
		return activitylog.NewService(logger), nil
	})

	app.Bind(application, "auth.service", func(application contracts.Application) (contracts.AuthService, error) {
		return auth.NewService(application), nil
	})

	app.Bind(application, "event.service", func(application contracts.Application) (contracts.EventService, error) {
		logger := app.MustMake[contracts.LoggingService](application, "logging.service")
		return event.NewService(application, logger), nil
	})

	app.Bind(application, "database.service", func(application contracts.Application) (contracts.DatabaseService, error) {
		return database.NewService(dbConfig)
	})

	app.Bind(application, "logging.service", func(application contracts.Application) (contracts.LoggingService, error) {
		return logging.NewGoLogger(log.Default()), nil
	})

	app.Bind(application, "crypt.service", func(application contracts.Application) (contracts.CryptService, error) {
		// key, err := base64.StdEncoding.DecodeString(cryptConfig.Key)
		// if err != nil {
		// 	return nil, err
		// }
		// if len([]byte(key)) != 32 {
		// 	return nil, fmt.Errorf("key length must be 32 bytes. generate key using `go run script/script.go key:generate`")
		// }

		return nil, nil
		// return crypt.NewAesGcmEncryptionService(key)
	})

	app.Bind(application, "http.middleware.service", func(a contracts.Application) (contracts.MiddlewareService, error) {
		return middleware.NewService(application, middlewareConfig), nil
	})

	app.Bind(application, "module.service", func(application contracts.Application) (contracts.ModuleService, error) {
		logger := app.MustMake[contracts.LoggingService](application, "logging.service")
		return module.NewService(application, logger), nil
	})

	app.Bind(application, "sessions.cookie_writer", func(application contracts.Application) (contracts.SessionCookieWriter, error) {
		return sessions.NewCookieUtil(sessionsConfig.Cookie), nil
	})
	app.Bind(application, "sessions.storage.database", func(a contracts.Application) (contracts.SessionStorage, error) {
		db := app.MustMake[contracts.DatabaseService](application, "database.service").GetDefault()
		return storage.NewDatabase(db, sessionsConfig.Table, sessionsConfig.AutoMigrate), nil
	})
	app.Bind(application, "sessions.storage.firestore", func(a contracts.Application) (contracts.SessionStorage, error) {
		client := app.MustMake[*firestore.Client](application, "firestore.client")
		return storage.NewFirestore(client, sessionsConfig.Table), nil
	})
	app.Bind(application, "sessions.service", func(application contracts.Application) (contracts.SessionService, error) {
		writer := app.MustMake[contracts.SessionCookieWriter](application, "sessions.cookie_writer")

		storageKey := fmt.Sprintf("sessions.storage.%s", sessionsConfig.Storage)
		storage, err := app.Make[contracts.SessionStorage](application, storageKey)
		if err != nil {
			return nil, fmt.Errorf("session service: failed to configure storage \"%s\": %w", sessionsConfig.Storage, err)
		}

		service, err := sessions.NewService(
			storage,
			writer,
			sessionsConfig,
		)

		return service, err
	})

	app.Bind(application, "web.engine", func(a contracts.Application) (*web.Engine, error) {
		loggingService := app.MustMake[contracts.LoggingService](application, "logging.service")
		engine, err := web.SetupEngine(loggingService, webConfig)
		if err != nil {
			return nil, err
		}

		return engine, nil
	})

	if err := registerAuthGuard(application); err != nil {
		return err
	}
	if err := registerMiddlewares(application); err != nil {
		return err
	}
	middlewareService, err := app.Make[contracts.MiddlewareService](application, "http.middleware.service")
	if err != nil {
		return err
	}
	middlewares := middlewareService.Global()
	engine := app.MustMake[*web.Engine](application, "web.engine")
	for _, m := range middlewares {
		engine.Use(m)
	}
	return nil
}
