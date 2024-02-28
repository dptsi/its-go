package sessions

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/web"
)

const sessionDataContextKey = "sessions.data"

type CookieConfig struct {
	Name           string
	CsrfCookieName string
	Path           string
	Domain         string
	Secure         bool

	/* Here you may specify the number of minutes that you wish the session
	to be allowed to remain idle before it expires. If you want them
	to immediately expire on the browser closing, set that option. */
	Lifetime int
}

type Config struct {
	// Default Session Storage
	// Supported: "database"
	Storage string

	// Database connection
	Connection string

	// Database table
	Table       string
	AutoMigrate bool

	Cookie CookieConfig
}

type Service struct {
	storage contracts.SessionStorage
	writer  contracts.SessionCookieWriter
	cfg     Config
}

func NewService(storage contracts.SessionStorage, writer contracts.SessionCookieWriter, cfg Config) (*Service, error) {
	return &Service{
		storage: storage,
		writer:  writer,
		cfg:     cfg,
	}, nil
}

func (s *Service) Get(ctx *web.Context, key string) (interface{}, error) {
	data, err := s.get(ctx)
	if err != nil {
		return nil, fmt.Errorf("session service: get: %w", err)
	}

	value, exists := data.Get(key)
	if !exists {
		return nil, nil
	}

	return value, nil
}

func (s *Service) Put(ctx *web.Context, key string, value interface{}) error {
	data, err := s.get(ctx)
	if err != nil {
		return fmt.Errorf("session service: put: %w", err)
	}

	data.Set(key, value)

	return s.updateToContextAndStorage(ctx, data)
}

func (s *Service) Delete(ctx *web.Context, key string) error {
	data, err := s.get(ctx)
	if err != nil {
		return fmt.Errorf("session service: delete: %w", err)
	}

	data.Delete(key)

	return s.updateToContextAndStorage(ctx, data)
}

func (s *Service) Clear(ctx *web.Context) error {
	data, err := s.get(ctx)
	if err != nil {
		return fmt.Errorf("session service: clear: %w", err)
	}

	data.Clear()

	return s.updateToContextAndStorage(ctx, data)
}

func (s *Service) Regenerate(ctx *web.Context) error {
	data, err := s.get(ctx)
	if err != nil {
		return fmt.Errorf("session service: regenerate: %w", err)
	}

	if err := s.storage.Delete(ctx, data.Id()); err != nil {
		return err
	}

	data.RegenerateId()
	if err := s.updateToContextAndStorage(ctx, data); err != nil {
		return err
	}

	s.writer.Write(ctx, data)
	return nil
}

func (s *Service) Invalidate(ctx *web.Context) error {
	data, err := s.get(ctx)
	if err != nil {
		return fmt.Errorf("session service: invalidate: %w", err)
	}

	if err := s.storage.Delete(ctx, data.Id()); err != nil {
		return err
	}

	data.RegenerateId()
	data.Clear()

	if err := s.updateToContextAndStorage(ctx, data); err != nil {
		return err
	}

	s.writer.Write(ctx, data)
	return nil
}

func (s *Service) RegenerateToken(ctx *web.Context) error {
	data, err := s.get(ctx)
	if err != nil {
		return fmt.Errorf("session service: regenerate token: %w", err)
	}

	data.RegenerateCSRFToken()

	if err := s.updateToContextAndStorage(ctx, data); err != nil {
		return err
	}
	s.writer.Write(ctx, data)
	return nil
}

func (s *Service) updateToContextAndStorage(ctx *web.Context, data *Data) error {
	if err := s.storage.Save(ctx, data); err != nil {
		return err
	}

	ctx.Set(sessionDataContextKey, data)
	return nil
}

func (s *Service) IsTokenMatch(ctx *web.Context, token string) (bool, error) {
	data, err := s.get(ctx)
	if err != nil {
		return false, fmt.Errorf("session service: is token match: %w", err)
	}

	return data.CSRFToken() == token, nil
}

func (s *Service) get(ctx *web.Context) (*Data, error) {
	data := s.getFromContext(ctx)

	if data == nil {
		return nil, fmt.Errorf("session data not available, do you forgot to execute Start()")
	}

	return data, nil
}

func (s *Service) Start(ctx *web.Context) error {
	data, err := s.getFromStorage(ctx)
	if err != nil {
		return err
	}

	if data == nil {
		data = NewEmptyData(int64(s.cfg.Cookie.Lifetime))
	} else {
		// Update cookie expiration
		data.expiredAt = getExpirationFromMaxAge(int64(s.cfg.Cookie.Lifetime))
	}

	if err := s.storage.Save(ctx, data); err != nil {
		return err
	}

	ctx.Set(sessionDataContextKey, data)
	s.writer.Write(ctx, data)
	return nil
}

func (s *Service) getFromStorage(ctx *web.Context) (*Data, error) {
	// Initialize session data
	sessionId, err := ctx.Cookie(s.cfg.Cookie.Name)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return nil, err
	}

	// Get session data from storage
	sessInterface, err := s.storage.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	sessionData, _ := sessInterface.(*Data)

	return sessionData, nil
}

func (s *Service) getFromContext(ctx *web.Context) *Data {
	dataIf, exists := ctx.Get(sessionDataContextKey)
	if !exists {
		return nil
	}
	data, ok := dataIf.(*Data)
	if !ok {
		return nil
	}

	return data
}
