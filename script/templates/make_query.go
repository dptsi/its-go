package templates

const ModuleQueryObject = `package queries

type {{.NamePascalCase}} struct {
}

type {{.NamePascalCase}}Query interface {
}
`
