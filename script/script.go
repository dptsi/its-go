package script

import "bitbucket.org/dptsi/its-go/script/command"

func LoadFrameworkScripts(s *Service) {
	s.AddCommand(&command.GenerateAppKey{})
}
