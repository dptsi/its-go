package command

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
	"github.com/stoewer/go-strcase"
)

type MakeValueObject struct{}

func (c *MakeValueObject) Key() string {
	return "make:domain:vo"
}

func (c *MakeValueObject) Name() string {
	return "Create new value object"
}

func (c *MakeValueObject) Description() string {
	return "Create new value object and its boilerplate"
}

func (c *MakeValueObject) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_value_object_name>", c.Key())
}

func (c *MakeValueObject) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no value object name provided")
	}
	modName := args[0]
	snakeCaseName := args[1]
	mod, err := newModule(modName)

	if err != nil {
		return fmt.Errorf("error detecting module: %w", err)
	}

	if err := mod.checkExist(); err != nil {
		return fmt.Errorf("error detecting module: %w", err)
	}

	if err := c.createValueObjectFile(mod, snakeCaseName); err != nil {
		return fmt.Errorf("error when creating %s.go: %w", snakeCaseName, err)
	}

	return nil
}

func (c *MakeValueObject) createValueObjectFile(mod *module, snakeCaseName string) error {
	path := mod.joinPath("internal/domain/valueobjects")

	path = filepath.Join(path, fmt.Sprintf("%s.go", snakeCaseName))
	valueObject, err := os.Create(path)
	if err != nil {
		return err
	}
	defer valueObject.Close()
	template, err := template.New("module_value_object").Parse(templates.ModuleValueObject)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(valueObject, data); err != nil {
		return err
	}

	return valueObject.Sync()
}
