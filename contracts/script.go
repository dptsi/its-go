package contracts

type ScriptCommand interface {
	Key() string
	Name() string
	Description() string
	Usage() string
	Handler(args []string) error
}
