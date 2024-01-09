package command

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
	"github.com/stoewer/go-strcase"
)

type MakeQuery struct{}

func (c *MakeQuery) Key() string {
	return "make:query"
}

func (c *MakeQuery) Name() string {
	return "Create new query object"
}

func (c *MakeQuery) Description() string {
	return "Create new query object and its boilerplate"
}

func (c *MakeQuery) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_query_name>", c.Key())
}

func (c *MakeQuery) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no query object name provided")
	}
	modName := args[0]
	snakeCaseName := args[1]
	mod, err := newModule(modName)

	if err != nil {
		return fmt.Errorf("error creating module: %w", err)
	}

	if err := mod.checkExist(); err != nil {
		return fmt.Errorf("error creating module: %w", err)
	}

	if err := c.createQueryObjectFile(mod, snakeCaseName); err != nil {
		return fmt.Errorf("error when creating %s.go: %w", snakeCaseName, err)
	}

	return nil
}

func (c *MakeQuery) createQueryObjectFile(mod *module, snakeCaseName string) error {
	path := mod.joinPath("internal/app/queries")

	path = filepath.Join(path, fmt.Sprintf("%s_query.go", snakeCaseName))
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s already exists", path)
	}
	queryObject, err := os.Create(path)
	if err != nil {
		return err
	}
	defer queryObject.Close()
	template, err := template.New("module_query_object").Parse(templates.ModuleQueryObject)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(queryObject, data); err != nil {
		return err
	}

	return queryObject.Sync()
}
