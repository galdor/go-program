package main

import (
	"fmt"
	"strings"

	"go.n16f.net/program"
)

func main() {
	var c *program.Command

	p := program.NewProgram("commands",
		"an example program with commands")

	p.AddFlag("", "flag-a", "a long flag")
	p.AddFlag("b", "", "a short flag")
	p.AddOption("c", "option-c", "value", "foo",
		"an option with both a short and long name")

	c = p.AddCommand("foo", "foo command", cmdFoo)
	c.AddFlag("d", "flag-d", "a command flag")
	c.AddArgument("arg-1", "the first argument")
	c.AddArgument("arg-2", "the second argument")
	c.AddTrailingArgument("arg-3", "all trailing arguments")

	c = p.AddCommand("bar", "bar command", cmdBar)
	c.AddOptionalArgument("arg-opt", "the optional argument")

	p.ParseCommandLine()
	p.Run()
}

func cmdFoo(p *program.Program) {
	p.Info("running command %q", p.CommandFullName())

	t := program.NewKeyValueTable()
	t.AddRow("flag-a", p.IsOptionSet("flag-a"))
	t.AddRow("b", p.IsOptionSet("b"))
	t.AddRow("option-c", p.OptionValue("option-c"))
	t.AddRow("flag-d", p.IsOptionSet("flag-d"))
	t.AddRow("arg-1", p.ArgumentValue("arg-1"))
	t.AddRow("arg-2", p.ArgumentValue("arg-2"))
	t.AddRow("arg-3", strings.Join(p.TrailingArgumentValues("arg-3"), " "))
	t.Print()
}

func cmdBar(p *program.Program) {
	p.Info("running command %q", p.CommandFullName())

	t := program.NewKeyValueTable()
	t.AddRow("flag-a", p.IsOptionSet("flag-a"))
	t.AddRow("b", p.IsOptionSet("b"))
	t.AddRow("option-c", p.OptionValue("option-c"))
	t.AddRow("arg-opt", p.ArgumentValue("arg-opt"))
	t.Print()

	fmt.Printf("flag-a: %v\n", p.IsOptionSet("flag-a"))
	fmt.Printf("b: %v\n", p.IsOptionSet("b"))
	fmt.Printf("option-c: %s\n", p.OptionValue("option-c"))

	fmt.Printf("arg-opt: %s\n", p.ArgumentValue("arg-opt"))
}
