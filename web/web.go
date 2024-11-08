package web

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/dptsi/its-go/app/errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Config struct {
	IsDebugMode bool
	Environment string
	Port        string
}

func SetupEngine(logger ErrorLogger, cfg Config) (*Engine, error) {
	if cfg.IsDebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "" {
				return name
			}
			if name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]; name != "" {
				return name
			}
			return ""
		})
	}

	r.Use(func(ctx *Context) {
		ctx.Set("request_id", uuid.NewString())
	})

	r.NoRoute(func(ctx *Context) {
		ctx.Error(errors.NewNotFound(errors.NotFoundParam{}))
		ctx.Abort()
	})
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(ctx *Context) {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, H{
			"code":    http.StatusMethodNotAllowed,
			"message": "method_not_allowed",
			"data":    nil,
		})
	})
	r.Use(gin.CustomRecovery(func(ctx *Context, err any) {
		requestId, exists := ctx.Get("request_id")
		data := map[string]interface{}{
			"error": "server unable to handle error",
		}
		if exists {
			data["request_id"] = requestId
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, H{
			"code":    statusCode[internalServerError],
			"message": "internal_server_error",
			"data":    data,
		})
	}))
	r.Use(globalErrorHandler(logger, cfg.IsDebugMode))

	return r, nil
}
