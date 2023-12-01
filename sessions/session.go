package sessions

import (
	"bitbucket.org/dptsi/base-go-libraries/contracts"
)

func Default(ctx contracts.WebFrameworkContext) *Data {
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
