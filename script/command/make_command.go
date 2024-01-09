package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
	"github.com/stoewer/go-strcase"
)

type MakeCommand struct{}

func (c *MakeCommand) Key() string {
	return "make:command"
}

func (c *MakeCommand) Name() string {
	return "Create new command handler"
}

func (c *MakeCommand) Description() string {
	return "Create new command handler and its boilerplate"
}

func (c *MakeCommand) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_command_name>", c.Key())
}

func (c *MakeCommand) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no command name provided")
	}
	modName := args[0]
	snakeCaseName := args[1]
	suffix := "_command"
	if strings.HasSuffix(snakeCaseName, suffix) {
		snakeCaseName = strings.TrimSuffix(snakeCaseName, suffix)
		fmt.Printf("tidak perlu pakai suffix %s ya rekan2\n", suffix)
		fmt.Printf("suffix yang ada akan ditambahkan otomatis\n")
		fmt.Printf("kedepannya cukup jalankan script dengan `%s %s %s` saja\n", c.Key(), modName, snakeCaseName)
	}
	mod, err := newModule(modName)

	if err != nil {
		return fmt.Errorf("error detecting module: %w", err)
	}

	if err := mod.checkExist(); err != nil {
		return fmt.Errorf("error detecting module: %w", err)
	}

	if err := c.createCommandFile(mod, snakeCaseName); err != nil {
		return fmt.Errorf("error when creating %s.go: %w", snakeCaseName, err)
	}

	return nil
}

func (c *MakeCommand) createCommandFile(mod *module, snakeCaseName string) error {
	path := mod.joinPath("internal/app/commands")

	path = filepath.Join(path, fmt.Sprintf("%s.go", snakeCaseName))
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s already exists", path)
	}
	command, err := os.Create(path)
	if err != nil {
		return err
	}
	defer command.Close()
	template, err := template.New("module_command").Parse(templates.ModuleCommand)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(command, data); err != nil {
		return err
	}

	return command.Sync()
}
