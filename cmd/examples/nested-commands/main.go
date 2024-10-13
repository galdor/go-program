package main

import (
	"strings"

	"go.n16f.net/program"
)

func main() {
	var c *program.Command

	p := program.NewProgram("nested-commands",
		"an example program with nested commands")

	p.AddFlag("a", "", "a short flag")

	c = p.AddCommand("foo create", "create a foo", cmdFooCreate)
	c.AddFlag("-n", "dry-run", "only pretend to create the foo")
	c.AddArgument("name", "the name of the foo")
	c.AddTrailingArgument("option", "a creation option")

	c = p.AddCommand("foo delete", "delete a foo", cmdFooDelete)
	c.AddArgument("name", "the name of the foo")

	c = p.AddCommand("bar", "bar command", cmdBar)
	c.AddOptionalArgument("arg-opt", "the optional argument")

	p.ParseCommandLine()
	p.Run()
}

func cmdFooCreate(p *program.Program) {
	p.Info("running command %q", p.CommandFullName())

	t := program.NewKeyValueTable()
	t.AddRow("a", p.IsOptionSet("a"))
	t.AddRow("dry-run", p.IsOptionSet("dry-run"))
	t.AddRow("name", p.ArgumentValue("name"))
	t.AddRow("options", strings.Join(p.TrailingArgumentValues("option"), " "))
	t.Print()
}

func cmdFooDelete(p *program.Program) {
	p.Info("running command %q", p.CommandFullName())

	t := program.NewKeyValueTable()
	t.AddRow("a", p.IsOptionSet("a"))
	t.AddRow("name", p.ArgumentValue("name"))
	t.Print()
}

func cmdBar(p *program.Program) {
	p.Info("running command %q", p.CommandFullName())

	t := program.NewKeyValueTable()
	t.AddRow("a", p.IsOptionSet("a"))
	t.AddRow("arg-opt", p.ArgumentValue("arg-opt"))
	t.Print()
}
