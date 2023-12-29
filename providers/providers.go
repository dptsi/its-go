package providers

import (
	"fmt"
	"log"

	"bitbucket.org/dptsi/go-framework/app"
	"bitbucket.org/dptsi/go-framework/auth"
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/database"
	"bitbucket.org/dptsi/go-framework/sessions"
	"bitbucket.org/dptsi/go-framework/sessions/storage"
	"bitbucket.org/dptsi/go-framework/web"
)

func LoadProviders(application *app.Application) error {
	config := application.Config()
	dbConfig, ok := config["database"].(database.Config)
	if !ok {
		return fmt.Errorf("database config is not available")
	}
	sessionsConfig, ok := config["sessions"].(sessions.Config)
	if !ok {
		return fmt.Errorf("sessions config is not available")
	}
	webConfig, ok := config["web"].(web.Config)
	if !ok {
		return fmt.Errorf("web config is not available")
	}

	log.Println("Registering database service...")
	app.Bind[contracts.DatabaseService](application, "database.service", func(application *app.Application) (contracts.DatabaseService, error) {
		return database.NewService(dbConfig)
	})
	log.Println("Database service registered!")

	log.Println("Registering sessions service...")
	app.Bind[contracts.SessionCookieWriter](application, "sessions.cookie_writer", func(application *app.Application) (contracts.SessionCookieWriter, error) {
		return sessions.NewCookieUtil(sessionsConfig.Cookie), nil
	})
	app.Bind[contracts.SessionStorage](application, "sessions.storage.db", func(a *app.Application) (contracts.SessionStorage, error) {
		db := app.MustMake[contracts.DatabaseService](application, "database.service").GetDefault()
		return storage.NewDatabase(db, sessionsConfig.Table), nil
	})
	app.Bind[contracts.SessionService](application, "sessions.service", func(application *app.Application) (contracts.SessionService, error) {
		writer := app.MustMake[contracts.SessionCookieWriter](application, "sessions.cookie_writer")

		storageKey := fmt.Sprintf("sessions.storage.%s", sessionsConfig.Storage)
		storage, err := app.Make[contracts.SessionStorage](application, storageKey)
		if err != nil {
			return nil, fmt.Errorf("session service: storage\"%s\" is not supported: %w", sessionsConfig.Storage, err)
		}

		service, err := sessions.NewService(
			storage,
			writer,
			sessionsConfig,
		)

		return service, err
	})
	log.Println("Session service registered!")

	log.Println("Registering authentication service...")
	app.Bind[contracts.AuthGuard](application, "auth.guard.sessions", func(application *app.Application) (contracts.AuthGuard, error) {
		return auth.NewSessionGuard(app.MustMake[contracts.SessionService](application, "sessions.service")), nil
	})
	app.Bind[contracts.AuthService](application, "auth.service", func(application *app.Application) (contracts.AuthService, error) {
		service := auth.NewService(auth.Config{
			Guards: map[string]auth.GuardsConfig{
				"sessions": {
					Driver: app.MustMake[contracts.AuthGuard](application, "auth.guard.sessions"),
				},
			},
		})

		return service, nil
	})
	log.Println("Authentication service registered!")

	log.Println("Registering web server...")
	app.Bind[*web.Engine](application, "web.engine", func(a *app.Application) (*web.Engine, error) {
		return web.SetupEngine(webConfig)
	})
	log.Println("Web server registered!")

	return nil
}
