package command

type Command interface {
	Init([]string) error
	Run() error
	Name() string
}
