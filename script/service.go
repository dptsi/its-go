package script

import (
	"fmt"
	"os"

	"bitbucket.org/dptsi/its-go/contracts"
)

type Service struct {
	commands map[string]contracts.ScriptCommand
}

func NewScriptService() *Service {
	return &Service{
		commands: make(map[string]contracts.ScriptCommand),
	}
}

func (s *Service) AddCommand(c contracts.ScriptCommand) error {
	if _, isExist := s.commands[c.Key()]; isExist {
		return fmt.Errorf("command key %s for %s already exist", c.Key(), c.Name())
	}
	s.commands[c.Key()] = c

	return nil
}

func (s *Service) Run() error {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "help" {
		fmt.Println("Available commands:")
		fmt.Println()

		for _, cmd := range s.commands {
			fmt.Printf("======== %s ========\n", cmd.Name())
			fmt.Printf("Description\t: %s\n", cmd.Name())
			fmt.Printf("Usage\t\t: %s\n", cmd.Usage())
			fmt.Println()
		}

		if len(s.commands) == 0 {
			fmt.Println("- No command available yet.")
		}

		return nil
	}

	command, exist := s.commands[args[0]]
	if !exist {
		return fmt.Errorf("command %s not found\nrun this script without arguments or with argument \"help\" to show help", args[0])
	}
	return command.Handler()(args[1:])
}
