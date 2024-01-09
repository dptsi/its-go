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

type MakeController struct{}

func (c *MakeController) Key() string {
	return "make:controller"
}

func (c *MakeController) Name() string {
	return "Create new controller handler"
}

func (c *MakeController) Description() string {
	return "Create new controller handler and its boilerplate"
}

func (c *MakeController) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_controller_name>", c.Key())
}

func (c *MakeController) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no controller name provided")
	}
	modName := args[0]
	snakeCaseName := args[1]
	mod, err := newModule(modName)
	suffix := "_controller"
	if strings.HasSuffix(snakeCaseName, suffix) {
		snakeCaseName = strings.TrimSuffix(snakeCaseName, suffix)
		fmt.Printf("tidak perlu pakai suffix %s ya rekan2\n", suffix)
		fmt.Printf("suffix yang ada akan ditambahkan otomatis\n")
		fmt.Printf("kedepannya cukup jalankan script dengan `%s %s %s` saja\n", c.Key(), modName, snakeCaseName)
	}

	if err != nil {
		return fmt.Errorf("error detecting module: %w", err)
	}

	if err := mod.checkExist(); err != nil {
		return fmt.Errorf("error detecting module: %w", err)
	}

	if err := c.createControllerFile(mod, snakeCaseName); err != nil {
		return fmt.Errorf("error when creating %s.go: %w", snakeCaseName, err)
	}

	fmt.Printf("controller %s berhasil dibuat pada modul %s!\n", snakeCaseName, modName)
	return nil
}

func (c *MakeController) createControllerFile(mod *module, snakeCaseName string) error {
	path := mod.joinPath("internal/presentation/controllers")

	path = filepath.Join(path, fmt.Sprintf("%s_controller.go", snakeCaseName))
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s already exists", path)
	}
	controller, err := os.Create(path)
	if err != nil {
		return err
	}
	defer controller.Close()
	template, err := template.New("module_controller").Parse(templates.ModuleController)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(controller, data); err != nil {
		return err
	}

	return controller.Sync()
}
