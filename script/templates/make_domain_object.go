package templates

var ModuleValueObject = `package valueobjects

type {{.NamePascalCase}} struct {
}

func New{{.NamePascalCase}}() ({{.NamePascalCase}}, error) {
	return {{.NamePascalCase}}{}, nil
}
`

var ModuleEntity = `package entities

import "github.com/dptsi/its-go/contracts"

type {{.NamePascalCase}} struct {
	events  []contracts.Event
	version int
}

func New{{.NamePascalCase}}(
	version int,
) (*{{.NamePascalCase}}, error) {
	return &{{.NamePascalCase}}{
		events:  make([]contracts.Event, 0),
		version: version,
	}, nil
}
`
