package templates

var ModuleValueObject = `package valueobjects

type {{.NamePascalCase}} struct {
}

func New{{.NamePascalCase}}() ({{.NamePascalCase}}, error) {
	return {{.NamePascalCase}}{}, nil
}
`
