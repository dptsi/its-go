package templates

var ModuleQueryObject = `package queries

type {{.NamePascalCase}} struct {
}

type {{.NamePascalCase}}Query interface {
}
`
