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

type MakeRepository struct{}

func (c *MakeRepository) Key() string {
	return "make:domain:repo"
}

func (c *MakeRepository) Name() string {
	return "Create new repository"
}

func (c *MakeRepository) Description() string {
	return "Create new repository and its boilerplate"
}

func (c *MakeRepository) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name> <snake_case_repository_name>", c.Key())
}

func (c *MakeRepository) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	if len(args) == 1 {
		return fmt.Errorf("no repository name provided")
	}
	modName := args[0]
	snakeCaseName := args[1]
	mod, err := newModule(modName)
	suffix := "_repository"
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

	if err := c.createRepositoryFile(mod, snakeCaseName); err != nil {
		return fmt.Errorf("error when creating %s.go: %w", snakeCaseName, err)
	}

	return nil
}

func (c *MakeRepository) createRepositoryFile(mod *module, snakeCaseName string) error {
	path := mod.joinPath("internal/domain/repositories")

	path = filepath.Join(path, fmt.Sprintf("%s_repository.go", snakeCaseName))
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file %s already exists", path)
	}
	repository, err := os.Create(path)
	if err != nil {
		return err
	}
	defer repository.Close()
	template, err := template.New("module_repository").Parse(templates.ModuleRepositoryFile)
	if err != nil {
		return err
	}

	type Data struct {
		NamePascalCase string
	}
	data := Data{
		NamePascalCase: strcase.UpperCamelCase(snakeCaseName),
	}

	if err := template.Execute(repository, data); err != nil {
		return err
	}

	return repository.Sync()
}
