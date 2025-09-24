package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/go-playground/validator/v10"

	commonErrors "github.com/dptsi/its-go/app/errors"
)

type ErrorLogger interface {
	Error(ctx context.Context, message string) error
}

func globalErrorHandler(logger ErrorLogger, isDebugMode bool) HandlerFunc {
	return func(ctx *Context) {
		ctx.Next()
		err := ctx.Errors.Last()
		if err == nil {
			return
		}
		requestId := ""
		reqIdInterface, exists := ctx.Get("request_id")
		if exists {
			if reqId, ok := reqIdInterface.(string); ok {
				requestId = reqId
			}
		}

		data := H{
			"request_id": requestId,
		}

		var validationErrors validator.ValidationErrors
		var badRequestError commonErrors.BadRequest
		var notFoundError commonErrors.NotFound
		var aggregateVersionMismatchError commonErrors.AggregateVersionMismatch
		var invariantError commonErrors.Invariant
		var forbiddenErr commonErrors.Forbidden
		var unauthorizedErr commonErrors.Unauthorized
		if errors.As(err, &validationErrors) {
			errorData := commonErrors.GetValidationErrors(validationErrors)
			data["errors"] = errorData
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error()))
			ctx.JSON(
				http.StatusBadRequest,
				H{
					"code":    statusCode[validationError],
					"message": validationError,
					"data":    data,
				},
			)
		} else if errors.As(err, &badRequestError) {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error()))
			for key, val := range badRequestError.Data() {
				data[key] = val
			}
			ctx.JSON(
				http.StatusBadRequest,
				H{
					"code":    badRequestError.Code(),
					"message": badRequestError.Message(),
					"data":    data,
				},
			)
		} else if errors.As(err, &invariantError) {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error()))
			ctx.JSON(
				http.StatusBadRequest,
				H{
					"code":    invariantError.Code(),
					"message": invariantError.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &notFoundError) {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 404; Error: %s\n", requestId, err.Error()))
			ctx.JSON(
				http.StatusNotFound,
				H{
					"code":    notFoundError.Code(),
					"message": notFoundError.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &aggregateVersionMismatchError) {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 409; Error: %s\n", requestId, err.Error()))
			ctx.JSON(
				http.StatusConflict,
				H{
					"code":    aggregateVersionMismatchError.Code(),
					"message": aggregateVersionMismatchError.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &forbiddenErr) {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 403; Error: %s\n", requestId, err.Error()))
			if (isDebugMode || !forbiddenErr.IsDetailRemovedInProd()) && forbiddenErr.Details() != "" {
				data["error"] = forbiddenErr.Details()
			}
			ctx.JSON(
				http.StatusForbidden,
				H{
					"code":    statusCode[forbiddenError],
					"message": forbiddenErr.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &unauthorizedErr) {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 401; Error: %s\n", requestId, err.Error()))
			ctx.JSON(
				http.StatusUnauthorized,
				H{
					"code":    statusCode[unauthorizedError],
					"message": unauthorizedError,
					"data":    data,
				},
			)
		} else {
			logger.Error(ctx, fmt.Sprintf("Request ID: %s; Status: 500; Error: %s\n", requestId, err.Error()))
			if isDebugMode {
				data["error"] = err.Error()
			}
			ctx.JSON(
				http.StatusInternalServerError,
				H{
					"code":    statusCode[internalServerError],
					"message": internalServerError,
					"data":    data,
				},
			)
		}

		// if sentry or sentrygin is disabled, this does nothing.
		if sentryhub := sentrygin.GetHubFromContext(ctx); sentryhub != nil {
			sentryhub.WithScope(func(scope *sentry.Scope) {
				scope.AddEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
					// so that error is grouped per path in sentry issue GUI
					event.Transaction = ctx.Request.URL.Path
					// default per sentry-go v0.35.1, if event.Transaction is empty,
					// error is grouped to this its-go module
					// because there's no stack trace available on the
					// error (err) instance, or pointed to code
					// that gives error if stack trace is available.
					return event
				})
				defer sentryhub.CaptureException(errors.Unwrap(err))
			})
		}
		// panic does not go to this middleware,
		// it goes straight to a recovery middleware,
		// which is a separate middleware.

		ctx.Abort()
	}
}
