package command

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
	"github.com/stoewer/go-strcase"
)

type MakeEntity struct{}

func (c *MakeEntity) Key() string {
	return "make:domain:entity"
}

func (c *MakeEntity) Name() string {
	return "Create new entity"
}

func (c *MakeEntity) Description() string {
	return "Create new entity and its boilerplate"
}

func (c *MakeEntity) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_entity_name>", c.Key())
}

func (c *MakeEntity) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no entity name provided")
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

	if err := c.createEntityFile(mod, snakeCaseName); err != nil {
		return fmt.Errorf("error when creating %s.go: %w", snakeCaseName, err)
	}

	fmt.Printf("entity %s berhasil dibuat pada modul %s!\n", snakeCaseName, modName)
	return nil
}

func (c *MakeEntity) createEntityFile(mod *module, snakeCaseName string) error {
	path := mod.joinPath("internal/domain/entities")

	path = filepath.Join(path, fmt.Sprintf("%s.go", snakeCaseName))
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s already exists", path)
	}
	entity, err := os.Create(path)
	if err != nil {
		return err
	}
	defer entity.Close()
	template, err := template.New("module_entity").Parse(templates.ModuleEntity)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(entity, data); err != nil {
		return err
	}

	return entity.Sync()
}
