// Package responses contains helper functions to create response
// that can be returned by the application
package responses

import (
	"fmt"
	"net/http"
)

type InfiniteScrollLinks struct {
	Next string `json:"next" example:"http://localhost:8080?cursor=2021-08-01&limit=10"`
}

type InfiniteScrollMeta struct {
	Total int `json:"total" example:"100"`
}

// InfiniteScrollResponse is a struct that contains the response of infinite scroll
//
// This response can be used to return the result of infinite scroll by using the following code:
//
//	var result InfiniteScrollResult[User]
//	ctx.JSON(result.GetInfiniteScrollResponse(urlConfig, limit))
type InfiniteScrollResponse[T any] struct {
	Code    int                 `json:"code" example:"123"`
	Message string              `json:"message"`
	Links   InfiniteScrollLinks `json:"links"`
	Meta    InfiniteScrollMeta  `json:"meta"`
	Data    []T                 `json:"data"`
}

// InfiniteScrollResult is a struct that contains the result of infinite scroll
// from the database using query object
type InfiniteScrollResult[T any] struct {
	NextCursor string
	Data       []T
	Total      int
}

// GetInfiniteScrollResponse is a function to create InfiniteScrollResponse
// from InfiniteScrollResult
func (r *InfiniteScrollResult[T]) GetInfiniteScrollResponse(urlConfig UrlConfig, limit int) InfiniteScrollResponse[T] {
	urlConfig.Query.Set("cursor", r.NextCursor)
	urlConfig.Query.Set("limit", fmt.Sprintf("%d", limit))
	return InfiniteScrollResponse[T]{
		Code:    http.StatusOK,
		Message: "success",
		Links: InfiniteScrollLinks{
			Next: fmt.Sprintf("%s?%s", urlConfig.FullUrl(), urlConfig.Query.Encode()),
		},
		Meta: InfiniteScrollMeta{
			Total: r.Total,
		},
		Data: r.Data,
	}
}

func (r *InfiniteScrollResult[T]) GetResponseByBaseUrl(baseUrl string, limit int) InfiniteScrollResponse[T] {
	return InfiniteScrollResponse[T]{
		Code:    http.StatusOK,
		Message: "success",
		Links: InfiniteScrollLinks{
			Next: fmt.Sprintf("%s?cursor=%s&limit=%d", baseUrl, r.NextCursor, limit),
		},
		Meta: InfiniteScrollMeta{
			Total: r.Total,
		},
		Data: r.Data,
	}
}
