package command

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

const (
	goModFileName = "go.mod"
)

type module struct {
	Base string
	Name string
	Path string
}

func newModule(name string) (*module, error) {
	path := fmt.Sprintf("modules/%s", name)

	base, err := getBaseModule()
	if err != nil {
		return nil, fmt.Errorf("error when reading module name from %s: %w", goModFileName, err)
	}
	mod := &module{
		Base: base,
		Name: name,
		Path: path,
	}

	return mod, nil
}

func (m module) joinPath(path string) string {
	return fmt.Sprintf("%s/%s", m.Path, path)
}

func (m module) createFolder() error {
	if err := os.Mkdir(m.Path, os.ModePerm); errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("module %s already exist", m.Name)
	}

	return nil
}

func (m module) checkExist() error {
	if _, err := os.Stat(m.Path); errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("module is not exist")
	}

	return nil
}

func getBaseModule() (string, error) {
	goMod, err := os.Open(goModFileName)
	if err != nil {
		return "", err
	}
	defer goMod.Close()

	var baseMod string
	fmt.Fscanf(goMod, "module %s", &baseMod)

	return baseMod, goMod.Sync()
}
