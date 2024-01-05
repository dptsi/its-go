package command

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dptsi/its-go/script/templates"
)

type MakeModule struct{}

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
	mod, err := newModule(name)

	if err != nil {
		return fmt.Errorf("error creating module: %w", err)
	}

	if err := mod.createFolder(); err != nil {
		return fmt.Errorf("error creating module: %w", err)
	}

	if err := c.createModuleEntrypoint(mod); err != nil {
		return fmt.Errorf("error when creating module entrypoint: %w", err)
	}

	if err := c.createModuleDependenciesFile(mod); err != nil {
		return fmt.Errorf("error when creating dependencies.go: %w", err)
	}

	if err := c.createModuleRoutesFile(mod); err != nil {
		return fmt.Errorf("error when creating routes.go: %w", err)
	}

	if err := c.createModuleAuthFile(mod); err != nil {
		return fmt.Errorf("error when creating auth.go: %w", err)
	}

	if err := c.createModuleEventsFile(mod); err != nil {
		return fmt.Errorf("error when creating events.go: %w", err)
	}

	if err := c.createModuleMiddlewaresFile(mod); err != nil {
		return fmt.Errorf("error when creating middlewares.go: %w", err)
	}

	if err := c.createModuleDirectories(mod); err != nil {
		return fmt.Errorf("error when creating module directories: %w", err)
	}

	return nil
}

func (c *MakeModule) createModuleEntrypoint(mod *module) error {
	fileName := fmt.Sprintf("%s.go", mod.Name)
	init, err := os.Create(mod.joinPath(fileName))
	if err != nil {
		return err
	}
	defer init.Close()
	template, err := template.New("module_entrypoint").Parse(templates.ModuleEntrypointTemplate)
	if err != nil {
		return err
	}
	if err := template.Execute(init, mod); err != nil {
		return err
	}

	return init.Sync()
}

func (c *MakeModule) createModuleDependenciesFile(mod *module) error {
	path := mod.joinPath("internal/app/providers")

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
	if err := template.Execute(deps, mod); err != nil {
		return err
	}

	return deps.Sync()
}

func (c *MakeModule) createModuleRoutesFile(mod *module) error {
	path := mod.joinPath("internal/presentation/routes")

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
	if err := template.Execute(routes, mod); err != nil {
		return err
	}

	return routes.Sync()
}

func (c *MakeModule) createModuleAuthFile(mod *module) error {
	path := mod.joinPath("internal/app/providers")

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
	if err := template.Execute(auth, mod); err != nil {
		return err
	}

	return auth.Sync()
}

func (c *MakeModule) createModuleEventsFile(mod *module) error {
	path := mod.joinPath("internal/app/providers")

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
	if err := template.Execute(events, mod); err != nil {
		return err
	}

	return events.Sync()
}

func (c *MakeModule) createModuleMiddlewaresFile(mod *module) error {
	path := mod.joinPath("internal/app/providers")

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
	if err := template.Execute(middlewares, mod); err != nil {
		return err
	}

	return middlewares.Sync()
}

func (c *MakeModule) createModuleDirectories(mod *module) error {
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
		path := mod.joinPath(relativePath)
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
