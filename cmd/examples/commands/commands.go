package main

import (
	"fmt"

	"go.n16f.net/program"
)

func main() {
	var cmd *program.Command

	p := program.NewProgram("no-command",
		"an example program without any command")

	p.AddFlag("", "flag-a", "a long flag")
	p.AddFlag("b", "", "a short flag")
	p.AddOption("c", "option-c", "value", "foo",
		"an option with both a short and long name")

	cmd = p.AddCommand("foo", "foo command", cmdFoo)
	cmd.AddFlag("d", "flag-d", "a command flag")
	cmd.AddArgument("arg-1", "the first argument")
	cmd.AddArgument("arg-2", "the second argument")
	cmd.AddTrailingArgument("arg-3", "all trailing arguments")

	cmd = p.AddCommand("bar", "bar command", cmdBar)
	cmd.AddOptionalArgument("arg-opt", "the optional argument")

	p.ParseCommandLine()
	p.Run()
}

func cmdFoo(p *program.Program) {
	p.Info("running command foo")

	fmt.Printf("flag-a: %v\n", p.IsOptionSet("flag-a"))
	fmt.Printf("b: %v\n", p.IsOptionSet("b"))
	fmt.Printf("option-c: %s\n", p.OptionValue("option-c"))
	fmt.Printf("flag-d: %v\n", p.IsOptionSet("flag-d"))

	fmt.Printf("arg-1: %s\n", p.ArgumentValue("arg-1"))
	fmt.Printf("arg-2: %s\n", p.ArgumentValue("arg-2"))
	fmt.Printf("arg-3:")
	for _, value := range p.TrailingArgumentValues("arg-3") {
		fmt.Printf(" %s", value)
	}
	fmt.Printf("\n")
}

func cmdBar(p *program.Program) {
	p.Info("running command bar")

	fmt.Printf("flag-a: %v\n", p.IsOptionSet("flag-a"))
	fmt.Printf("b: %v\n", p.IsOptionSet("b"))
	fmt.Printf("option-c: %s\n", p.OptionValue("option-c"))

	fmt.Printf("arg-opt: %s\n", p.ArgumentValue("arg-opt"))
}
