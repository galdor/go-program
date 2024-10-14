package program

import (
	"maps"
	"os"
	"strings"
)

func (p *Program) parse() {
	args := p.parseOptions(os.Args[1:], p.options)

	if p.IsOptionSet("help") {
		return
	}

	if p.command == nil {
		args = p.parseArguments(args, p.arguments)
	} else {
		args = p.parseCommand(args)

		options := make(map[string]*Option)
		maps.Copy(options, p.options)
		maps.Copy(options, p.selectedCommand.options)

		args = p.parseOptions(args, options)

		args = p.parseArguments(args, p.selectedCommand.arguments)
	}
}

func (p *Program) parseOptions(args []string, options map[string]*Option) []string {
	for len(args) > 0 {
		arg := args[0]
		if arg == "--" || !isOption(arg) {
			break
		}

		key := strings.TrimLeft(arg, "-")

		opt, found := options[key]
		if !found {
			p.Fatal("unknown option %q", key)
		}

		opt.Set = true

		if opt.ValueName == "" {
			args = args[1:]
		} else {
			if len(args) < 2 {
				p.Fatal("missing value for option %q", key)
			}

			opt.Value = args[1]

			args = args[2:]
		}
	}

	return args
}

func (p *Program) parseCommand(args []string) []string {
	p.selectedCommand = p.command

	if len(args) == 0 {
		p.Fatal("missing command")
	}

	cmd := p.command
	names := []string{}

	for len(args) > 0 {
		arg := args[0]
		if arg == "--" || isOption(arg) {
			break
		}

		names = append(names, arg)

		cmd2 := cmd.subcommands[arg]
		if cmd2 == nil {
			break
		}

		cmd = cmd2
		args = args[1:]
	}

	fullName := strings.Join(names, " ")

	if cmd == p.command {
		p.Fatal("unknown command %q", fullName)
	}

	if cmd.Main == nil {
		if len(args) == 0 {
			p.Fatal("missing subcommand(s) for command %q", cmd.FullName)
		} else {
			p.Fatal("unknown command %q", strings.Join(names, " "))
		}
	}

	p.selectedCommand = cmd

	return args
}

func (p *Program) parseArguments(args []string, arguments []*Argument) []string {
	if len(arguments) > 0 {
		// Mandatory arguments
		min := 0
		for _, argument := range arguments {
			if argument.Optional || argument.Trailing {
				break
			}

			min++
		}

		if len(args) < min {
			p.Fatal("missing argument(s)")
		}

		for i := 0; i < min; i++ {
			argument := arguments[i]

			argument.Set = true
			argument.Value = args[i]
		}

		args = args[min:]
		arguments = arguments[min:]

		// Optional arguments
		var trailingArgument *Argument

		for _, argument := range arguments {
			if len(args) == 0 {
				break
			}

			if argument.Trailing {
				trailingArgument = argument
				break
			}

			argument.Set = true
			argument.Value = args[0]

			args = args[1:]
		}

		// Trailing argument
		if trailingArgument != nil {
			trailingArgument.TrailingValues = args
			args = args[len(args):]
		} else {
			if len(args) > 0 {
				p.Fatal("too many arguments")
			}
		}
	} else {
		if len(args) > 0 {
			p.Fatal("unexpected arguments")
		}
	}

	return args
}

func isOption(arg string) bool {
	return isShortOption(arg) || isLongOption(arg)
}

func isShortOption(arg string) bool {
	return len(arg) == 2 && arg[0] == '-' && arg[1] != '-'
}

func isLongOption(arg string) bool {
	return len(arg) > 2 && arg[0] == '-' && arg[1] == '-'
}
