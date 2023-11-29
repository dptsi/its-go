package middleware

import (
	"errors"
	"fmt"

	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"github.com/gin-gonic/gin"
)

func StartSession(storage sessions.Storage, attr sessions.AddSessionCookieToResponseAttributes) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if storage == nil {
			err := errors.New("session storage not configured. please configure it first in bootstrap/web/web.go")
			ctx.Error(fmt.Errorf("start session middleware: %w", err))
			ctx.Abort()
		}

		// Initialize session data
		var data *sessions.Data
		sessionId, err := ctx.Cookie(attr.Name)

		if err == nil {
			// Get session data from storage
			sess, err := storage.Get(ctx, sessionId)
			if err != nil {
				ctx.Error(err)
				ctx.Abort()
				return
			}
			if sess != nil {
				data = sess
			}
		}
		if data == nil {
			data = sessions.NewEmptyData(int64(attr.MaxAge))
			if err := storage.Save(ctx, data); err != nil {
				ctx.Error(fmt.Errorf("start session middleware: %w", err))
				ctx.Abort()
			}
		}
		ctx.Set("session", data)
		sessions.AddSessionCookieToResponse(ctx, attr, data)
		ctx.Next()
	}
}
