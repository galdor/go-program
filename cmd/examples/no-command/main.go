package main

import (
	"strings"

	"go.n16f.net/program"
)

func main() {
	p := program.NewProgram("no-command",
		"an example program without any command")

	p.AddFlag("", "flag-a", "a long flag")
	p.AddFlag("b", "", "a short flag")
	p.AddOption("c", "option-c", "value", "foo",
		"an option with both a short and long name")

	p.AddArgument("arg-1", "the first argument")
	p.AddArgument("arg-2", "the second argument")
	p.AddOptionalArgument("arg-opt-1", "the first optional argument")
	p.AddOptionalArgument("arg-opt-2", "the second optional argument")
	p.AddTrailingArgument("arg-trailing", "all trailing arguments")

	p.SetMain(main2)

	p.ParseCommandLine()

	p.Debug(2, "running program")

	p.Run()
}

func main2(p *program.Program) {
	t := program.NewKeyValueTable()
	t.AddRow("flag-a", p.IsOptionSet("flag-a"))
	t.AddRow("b", p.IsOptionSet("b"))
	t.AddRow("option-c", p.OptionValue("option-c"))
	t.AddRow("arg-1", p.ArgumentValue("arg-1"))
	t.AddRow("arg-2", p.ArgumentValue("arg-2"))
	t.AddRow("arg-opt-1", p.ArgumentValue("arg-opt-1"))
	t.AddRow("arg-opt-2", p.ArgumentValue("arg-opt-2"))
	t.AddRow("arg-trailing",
		strings.Join(p.TrailingArgumentValues("arg-trailing"), " "))
	t.Print()
}
