package web

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	commonErrors "github.com/dptsi/its-go/app/errors"
)

func globalErrorHandler(isDebugMode bool) HandlerFunc {
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
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusBadRequest,
				H{
					"code":    statusCode[validationError],
					"message": validationError,
					"data":    data,
				},
			)
		} else if errors.As(err, &badRequestError) {
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
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
			log.Printf("Request ID: %s; Status: 400; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusBadRequest,
				H{
					"code":    invariantError.Code(),
					"message": invariantError.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &notFoundError) {
			log.Printf("Request ID: %s; Status: 404; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusNotFound,
				H{
					"code":    notFoundError.Code(),
					"message": notFoundError.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &aggregateVersionMismatchError) {
			log.Printf("Request ID: %s; Status: 409; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusConflict,
				H{
					"code":    aggregateVersionMismatchError.Code(),
					"message": aggregateVersionMismatchError.Error(),
					"data":    data,
				},
			)
		} else if errors.As(err, &forbiddenErr) {
			log.Printf("Request ID: %s; Status: 403; Error: %s\n", requestId, err.Error())
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
			log.Printf("Request ID: %s; Status: 401; Error: %s\n", requestId, err.Error())
			ctx.JSON(
				http.StatusUnauthorized,
				H{
					"code":    statusCode[unauthorizedError],
					"message": unauthorizedError,
					"data":    data,
				},
			)
		} else {
			log.Printf("Request ID: %s; Status: 500; Error: %s\n", requestId, err.Error())
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

		ctx.Abort()
	}
}
