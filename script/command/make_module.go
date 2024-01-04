package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
)

const (
	goModFileName = "go.mod"
)

type MakeModule struct{}

type moduleConfig struct {
	BaseMod string
	Name    string
	Path    string
}

func (m moduleConfig) joinPath(path string) string {
	return fmt.Sprintf("%s/%s", m.Path, path)
}

func (c *MakeModule) Key() string {
	return "make:module"
}

func (c *MakeModule) Name() string {
	return "Create new module"
}

func (c *MakeModule) Description() string {
	return "Create new module and its boilerplate"
}

func (c *MakeModule) Usage() string {
	return fmt.Sprintf("%s <snake_case_module_name>", c.Key())
}

func (c *MakeModule) Handler(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no module name provided")
	}
	name := args[0]
	path := fmt.Sprintf("modules/%s", name)

	if err := os.Mkdir(path, os.ModePerm); errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("module %s already exist", name)
	}
	baseMod, err := c.getBaseModule()
	if err != nil {
		return fmt.Errorf("error when reading module name from %s: %w", goModFileName, err)
	}
	cfg := moduleConfig{
		BaseMod: baseMod,
		Name:    name,
		Path:    path,
	}

	if err := c.createModuleEntrypoint(cfg); err != nil {
		return fmt.Errorf("error when creating module entrypoint: %w", err)
	}

	if err := c.createModuleDependenciesFile(cfg); err != nil {
		return fmt.Errorf("error when creating dependencies.go: %w", err)
	}

	if err := c.createModuleRoutesFile(cfg); err != nil {
		return fmt.Errorf("error when creating routes.go: %w", err)
	}

	if err := c.createModuleAuthFile(cfg); err != nil {
		return fmt.Errorf("error when creating auth.go: %w", err)
	}

	if err := c.createModuleEventsFile(cfg); err != nil {
		return fmt.Errorf("error when creating events.go: %w", err)
	}

	if err := c.createModuleMiddlewaresFile(cfg); err != nil {
		return fmt.Errorf("error when creating middlewares.go: %w", err)
	}

	if err := c.createModuleDirectories(cfg); err != nil {
		return fmt.Errorf("error when creating module directories: %w", err)
	}

	return nil
}

func (c *MakeModule) getBaseModule() (string, error) {
	goMod, err := os.Open(goModFileName)
	if err != nil {
		return "", err
	}
	defer goMod.Close()

	var baseMod string
	fmt.Fscanf(goMod, "module %s", &baseMod)

	return baseMod, goMod.Sync()
}

func (c *MakeModule) createModuleEntrypoint(cfg moduleConfig) error {
	fileName := fmt.Sprintf("%s.go", cfg.Name)
	init, err := os.Create(cfg.joinPath(fileName))
	if err != nil {
		return err
	}
	defer init.Close()
	template, err := template.New("module_entrypoint").Parse(templates.ModuleEntrypointTemplate)
	if err != nil {
		return err
	}
	if err := template.Execute(init, cfg); err != nil {
		return err
	}

	return init.Sync()
}

func (c *MakeModule) createModuleDependenciesFile(cfg moduleConfig) error {
	path := cfg.joinPath("internal/app/providers")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	path = filepath.Join(path, "dependencies.go")
	deps, err := os.Create(path)
	if err != nil {
		return err
	}
	defer deps.Close()
	template, err := template.New("module_deps").Parse(templates.ModuleDeps)
	if err != nil {
		return err
	}
	if err := template.Execute(deps, cfg); err != nil {
		return err
	}

	return deps.Sync()
}

func (c *MakeModule) createModuleRoutesFile(cfg moduleConfig) error {
	path := cfg.joinPath("internal/presentation/routes")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	path = filepath.Join(path, "routes.go")
	routes, err := os.Create(path)
	if err != nil {
		return err
	}
	defer routes.Close()
	template, err := template.New("module_routes").Parse(templates.ModuleRoutes)
	if err != nil {
		return err
	}
	if err := template.Execute(routes, cfg); err != nil {
		return err
	}

	return routes.Sync()
}

func (c *MakeModule) createModuleAuthFile(cfg moduleConfig) error {
	path := cfg.joinPath("internal/app/providers")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	path = filepath.Join(path, "auth.go")
	auth, err := os.Create(path)
	if err != nil {
		return err
	}
	defer auth.Close()
	template, err := template.New("module_auth").Parse(templates.ModuleAuth)
	if err != nil {
		return err
	}
	if err := template.Execute(auth, cfg); err != nil {
		return err
	}

	return auth.Sync()
}

func (c *MakeModule) createModuleEventsFile(cfg moduleConfig) error {
	path := cfg.joinPath("internal/app/providers")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	path = filepath.Join(path, "events.go")
	events, err := os.Create(path)
	if err != nil {
		return err
	}
	defer events.Close()
	template, err := template.New("module_events").Parse(templates.ModuleEvent)
	if err != nil {
		return err
	}
	if err := template.Execute(events, cfg); err != nil {
		return err
	}

	return events.Sync()
}

func (c *MakeModule) createModuleMiddlewaresFile(cfg moduleConfig) error {
	path := cfg.joinPath("internal/app/providers")

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	path = filepath.Join(path, "middlewares.go")
	middlewares, err := os.Create(path)
	if err != nil {
		return err
	}
	defer middlewares.Close()
	template, err := template.New("module_middlewares").Parse(templates.ModuleMiddleware)
	if err != nil {
		return err
	}
	if err := template.Execute(middlewares, cfg); err != nil {
		return err
	}

	return middlewares.Sync()
}

func (c *MakeModule) createModuleDirectories(cfg moduleConfig) error {
	var moduleFolders = []string{
		"internal/app/commands",
		"internal/app/listeners",
		"internal/app/queries",
		"internal/app/services",

		"internal/infrastructures",

		"internal/domain/entities",
		"internal/domain/events",
		"internal/domain/valueobjects",
		"internal/domain/services",
		"internal/domain/repositories",

		"internal/presentation/controllers",
	}

	for _, relativePath := range moduleFolders {
		path := cfg.joinPath(relativePath)
		os.MkdirAll(path, os.ModePerm)
		gitKeepPath := filepath.Join(path, ".gitkeep")
		gitKeep, err := os.Create(gitKeepPath)
		if err != nil {
			return err
		}
		if err := gitKeep.Sync(); err != nil {
			return err
		}

		if err := gitKeep.Close(); err != nil {
			return err
		}
	}

	return nil
}
