package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

const DefaultRedisConnectionName = "default"

type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DB       int
}

type Config struct {
	Connections map[string]ConnectionConfig
}

type Service struct {
	clients map[string]*Client
}

func NewService(cfg Config) (*Service, error) {
	clients := make(map[string]*Client, len(cfg.Connections))

	for name, conn := range cfg.Connections {
		if name == "" {
			name = DefaultRedisConnectionName
		}
		if conn.Host == "" {
			continue
		}

		port := conn.Port
		if port == "" {
			port = "6379"
		}

		clients[name] = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", conn.Host, port),
			Username: conn.Username,
			Password: conn.Password,
			DB:       conn.DB,
		})
	}

	return &Service{clients: clients}, nil
}

func (s *Service) GetClient(name string) *Client {
	if s.clients == nil {
		return nil
	}
	return s.clients[name]
}

func (s *Service) GetDefault() *Client {
	return s.GetClient(DefaultRedisConnectionName)
}

func (s *Service) HasClient(name string) bool {
	return s.GetClient(name) != nil
}

func (s *Service) Close() error {
	if s.clients == nil {
		return nil
	}
	for _, client := range s.clients {
		if client != nil {
			_ = client.Close()
		}
	}
	return nil
}
