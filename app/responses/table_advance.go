package responses

import (
	"fmt"
	"net/http"
)

type TableAdvancedLinks struct {
	Next string `json:"next" example:"http://localhost:8080?start_after=2021-08-01&page=3"`
	Prev string `json:"prev" example:"http://localhost:8080?end_before=2021-08-01&page=1"`
}

type TableAdvancedMeta struct {
	Total int `json:"total" example:"100"`
	Range int `json:"range" example:"10"`
	Page  int `json:"page" example:"2"`
}

// TableAdvancedResponse is a struct that contains the response of table advanced
//
// This response can be used to return the result of table advanced by using the following code:
//
//	var result TableAdvancedResult[User]
//	ctx.JSON(result.GetTableAdvancedResponse(urlConfig, limit, currentPage))
type TableAdvancedResponse[T any] struct {
	Code    int                `json:"code" example:"123"`
	Message string             `json:"message"`
	Links   TableAdvancedLinks `json:"links"`
	Meta    TableAdvancedMeta  `json:"meta"`
	Data    []T                `json:"data"`
}

// TableAdvancedResult is a struct that contains the result of table advanced
// from the database using query object
type TableAdvancedResult[T any] struct {
	StartAfter string
	EndBefore  string
	Data       []T
	Total      int
	ItemCount  int
}

func (r *TableAdvancedResult[T]) GetTableAdvancedResponse(urlConfig UrlConfig, limit int, currentPage int) TableAdvancedResponse[T] {
	return TableAdvancedResponse[T]{
		Code:    http.StatusOK,
		Message: "success",
		Links: TableAdvancedLinks{
			Prev: fmt.Sprintf("%s?end_before=%s&page=%d", urlConfig.FullUrl(), r.EndBefore, max(1, currentPage-1)),
			Next: fmt.Sprintf("%s?start_after=%s&page=%d", urlConfig.FullUrl(), r.StartAfter, min(r.Total/limit+1, currentPage+1)),
		},
		Meta: TableAdvancedMeta{
			Total: r.Total,
			Range: r.ItemCount,
			Page:  currentPage,
		},
	}
}
