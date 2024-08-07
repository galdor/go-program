package main

import (
	"fmt"

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
	fmt.Printf("flag-a: %v\n", p.IsOptionSet("flag-a"))
	fmt.Printf("b: %v\n", p.IsOptionSet("b"))
	fmt.Printf("option-c: %s\n", p.OptionValue("option-c"))

	fmt.Printf("arg-1: %s\n", p.ArgumentValue("arg-1"))
	fmt.Printf("arg-2: %s\n", p.ArgumentValue("arg-2"))

	fmt.Printf("arg-opt-1: %s\n", p.ArgumentValue("arg-opt-1"))
	fmt.Printf("arg-opt-2: %s\n", p.ArgumentValue("arg-opt-2"))

	fmt.Printf("arg-trailing:")
	for _, value := range p.TrailingArgumentValues("arg-trailing") {
		fmt.Printf(" %s", value)
	}
	fmt.Printf("\n")
}
