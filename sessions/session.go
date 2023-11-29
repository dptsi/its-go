package sessions

import "github.com/gin-gonic/gin"

func Default(ctx *gin.Context) *Data {
	dataIf, exists := ctx.Get("session")
	if !exists {
		panic("session not found in context, make sure you have called session.StartSession middleware")
	}
	data, ok := dataIf.(*Data)
	if !ok {
		panic("session is not of type session.Data")
	}

	return data
}
