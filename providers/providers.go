package providers

import (
	"fmt"
	"log"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/auth"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/database"
	"bitbucket.org/dptsi/go-framework/event"
	"bitbucket.org/dptsi/go-framework/http/middleware"
	"bitbucket.org/dptsi/go-framework/module"
	"bitbucket.org/dptsi/go-framework/sessions"
	"bitbucket.org/dptsi/go-framework/sessions/storage"
	"bitbucket.org/dptsi/go-framework/web"
)

func LoadProviders(application contracts.Application) error {
	config := application.Config()
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

	log.Println("Registering authentication service...")
	app.Bind[contracts.AuthService](application, "auth.service", func(application contracts.Application) (contracts.AuthService, error) {
		return auth.NewService(application), nil
	})
	log.Println("Authentication service registered!")

	log.Println("Registering event service...")
	app.Bind[contracts.EventService](application, "event.service", func(application contracts.Application) (contracts.EventService, error) {
		return event.NewService(application), nil
	})
	log.Println("Event service registered!")

	log.Println("Registering database service...")
	app.Bind[contracts.DatabaseService](application, "database.service", func(application contracts.Application) (contracts.DatabaseService, error) {
		return database.NewService(dbConfig)
	})
	log.Println("Database service registered!")

	log.Println("Registering middleware service...")
	app.Bind[contracts.MiddlewareService](application, "http.middleware.service", func(a contracts.Application) (contracts.MiddlewareService, error) {
		return middleware.NewService(application, middlewareConfig), nil
	})
	log.Println("Middleware service registered!")

	log.Println("Registering module service...")
	app.Bind[contracts.ModuleService](application, "module.service", func(application contracts.Application) (contracts.ModuleService, error) {
		return module.NewService(application), nil
	})
	log.Println("Module service registered!")

	log.Println("Registering sessions service...")
	app.Bind[contracts.SessionCookieWriter](application, "sessions.cookie_writer", func(application contracts.Application) (contracts.SessionCookieWriter, error) {
		return sessions.NewCookieUtil(sessionsConfig.Cookie), nil
	})
	app.Bind[contracts.SessionStorage](application, "sessions.storage.database", func(a contracts.Application) (contracts.SessionStorage, error) {
		db := app.MustMake[contracts.DatabaseService](application, "database.service").GetDefault()
		return storage.NewDatabase(db, sessionsConfig.Table, sessionsConfig.AutoMigrate), nil
	})
	app.Bind[contracts.SessionService](application, "sessions.service", func(application contracts.Application) (contracts.SessionService, error) {
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
	log.Println("Session service registered!")

	log.Println("Registering web server...")
	app.Bind[*web.Engine](application, "web.engine", func(a contracts.Application) (*web.Engine, error) {
		engine, err := web.SetupEngine(webConfig)
		if err != nil {
			return nil, err
		}

		return engine, nil
	})
	log.Println("Web server registered!")

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
