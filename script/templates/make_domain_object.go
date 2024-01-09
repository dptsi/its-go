package templates

import "fmt"

const ModuleValueObject = `package valueobjects

type {{.NamePascalCase}} struct {
}

func New{{.NamePascalCase}}() ({{.NamePascalCase}}, error) {
	return {{.NamePascalCase}}{}, nil
}
`

const ModuleEntity = `package entities

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

var ModuleEventFile = fmt.Sprintf(`package events

import (
	"encoding/json"
	"time"
)

type {{.NamePascalCase}} struct {
	Timestamp time.Time %s
}

func New{{.NamePascalCase}}() {{.NamePascalCase}} {
	return {{.NamePascalCase}}{
		Timestamp: time.Now(),
	}
}

func (p {{.NamePascalCase}}) OccuredOn() time.Time {
	return p.Timestamp
}

func (p {{.NamePascalCase}}) JSON() ([]byte, error) {
	return json.Marshal(p)
}
`,
	"`json:\"timestamp\"`",
)
