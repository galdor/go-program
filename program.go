package program

import (
	"fmt"
	"os"
)

type Main func(*Program)

type Program struct {
	Name        string
	Description string
	Main        Main

	commands  map[string]*Command
	options   map[string]*Option
	arguments []*Argument

	command *Command

	Quiet      bool
	DebugLevel int
}

func NewProgram(name, description string) *Program {
	p := &Program{
		Name:        name,
		Description: description,

		commands: make(map[string]*Command),

		options: make(map[string]*Option),
	}

	p.addDefaultOptions()

	return p
}

func (p *Program) SetMain(main Main) {
	if len(p.commands) > 0 {
		panic("cannot have a main function with commands")
	}

	p.Main = main
}

func (p *Program) Run() {
	var main Main
	if p.command == nil {
		if p.Main == nil {
			panic("missing main function")
		}

		main = p.Main
	} else {
		main = p.command.Main
	}

	main(p)
}

func (p *Program) Debug(level int, format string, args ...interface{}) {
	if level > p.DebugLevel {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (p *Program) Info(format string, args ...interface{}) {
	if p.Quiet {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (p *Program) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
}

func (p *Program) Fatal(format string, args ...interface{}) {
	p.Error(format, args...)
	os.Exit(1)
}

func (p *Program) fatal(format string, args ...interface{}) {
	p.Error(format, args...)

	fmt.Fprintf(os.Stderr, "\n")

	if p.command == nil {
		p.PrintUsage(nil)
	} else {
		p.PrintUsage(p.command)
	}

	os.Exit(1)
}
