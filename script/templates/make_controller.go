package templates

const ModuleController = `package controllers

import (
	"net/http"

	"github.com/dptsi/its-go/web"
)

type {{.NamePascalCase}}Controller struct {
	// Tambahkan dependency yang diperlukan disini
}

func New{{.NamePascalCase}}Controller(
// Tambahkan dependency yang diperlukan disini
) *{{.NamePascalCase}}Controller {
	return &{{.NamePascalCase}}Controller{
		// Tambahkan dependency yang diperlukan disini
	}
}

func (c *{{.NamePascalCase}}Controller) Hello(ctx *web.Context) {
	ctx.JSON(http.StatusOK, web.H{
		"code":    1,
		"message": "hello",
		"data":    nil,
	})
}
`
