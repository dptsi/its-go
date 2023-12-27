package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HandleCors struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func (h *HandleCors) Execute(ctx *gin.Context) {
	cors := cors.New(cors.Config{
		AllowOrigins:     h.AllowedOrigins,
		AllowMethods:     h.AllowedMethods,
		AllowHeaders:     h.AllowedHeaders,
		ExposeHeaders:    h.ExposedHeaders,
		AllowCredentials: h.AllowCredentials,
		MaxAge:           time.Duration(h.MaxAge) * time.Second,
	})

	cors(ctx)
}
