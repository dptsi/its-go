package command

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
	"github.com/stoewer/go-strcase"
)

type MakeEvent struct{}

func (c *MakeEvent) Key() string {
	return "make:domain:event"
}

func (c *MakeEvent) Name() string {
	return "Create new event"
}

func (c *MakeEvent) Description() string {
	return "Create new event and its boilerplate"
}

func (c *MakeEvent) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_event_name>", c.Key())
}

func (c *MakeEvent) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no event name provided")
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

	path := filepath.Join(mod.joinPath("internal/domain/events"), fmt.Sprintf("%s.go", snakeCaseName))
	if err := c.createEventFile(mod, snakeCaseName, path); err != nil {
		return fmt.Errorf("error when creating %s: %w", path, err)
	}

	fmt.Printf("event %s berhasil dibuat pada %s!\n", snakeCaseName, path)
	return nil
}

func (c *MakeEvent) createEventFile(mod *module, snakeCaseName, path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s already exists", path)
	}
	event, err := os.Create(path)
	if err != nil {
		return err
	}
	defer event.Close()
	template, err := template.New("module_event").Parse(templates.ModuleEventFile)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(event, data); err != nil {
		return err
	}

	return event.Sync()
}
