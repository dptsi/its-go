package contracts

type ScriptCommandHandler = func(args []string) error

type ScriptCommand interface {
	Key() string
	Name() string
	Description() string
	Usage() string
	Handler() ScriptCommandHandler
}
