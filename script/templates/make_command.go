package templates

const ModuleCommand = `package commands

import "context"

type {{.NamePascalCase}}Request struct {
	Ctx context.Context
	// Tambahkan data lainnya disini
}

type {{.NamePascalCase}}Command struct {
	// Tambahkan dependency yang diperlukan disini
}

func New{{.NamePascalCase}}Command() *{{.NamePascalCase}}Command {
	return &{{.NamePascalCase}}Command{}
}

// Modifikasi return value dari Execute() seperlunya
// Misalkan Anda ingin mengembalikan id dari entity
// cukup ubah menjadi Execute(req {{.NamePascalCase}}Request) (id string, err error)
func (c *{{.NamePascalCase}}Command) Execute(req {{.NamePascalCase}}Request) (err error) {
	return nil
}
`
