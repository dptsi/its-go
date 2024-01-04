package middleware

import (
	"time"

	"github.com/dptsi/its-go/http"
	"github.com/dptsi/its-go/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Cors struct {
	cfg http.CorsConfig
}

func NewCors(cfg http.CorsConfig) *Cors {
	return &Cors{cfg}
}

func (h *Cors) Handle(interface{}) web.HandlerFunc {
	return func(ctx *gin.Context) {
		cors := cors.New(cors.Config{
			AllowOrigins:     h.cfg.AllowedOrigins,
			AllowMethods:     h.cfg.AllowedMethods,
			AllowHeaders:     h.cfg.AllowedHeaders,
			ExposeHeaders:    h.cfg.ExposedHeaders,
			AllowCredentials: h.cfg.AllowCredentials,
			MaxAge:           time.Duration(h.cfg.MaxAge) * time.Second,
		})

		cors(ctx)
	}
}
