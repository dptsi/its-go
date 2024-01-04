package script

import "github.com/dptsi/its-go/script/command"

func LoadFrameworkScripts(s *Service) {
	s.AddCommand(&command.GenerateAppKey{})
	s.AddCommand(&command.MakeModule{})
}
