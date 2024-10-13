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

		isShort := len(arg) == 2 && arg[0] == '-' && arg[1] != '-'
		isLong := len(arg) > 2 && arg[0:2] == "--"

		if arg == "--" || !(isShort || isLong) {
			break
		}

		key := strings.TrimLeft(arg, "-")

		opt, found := options[key]
		if !found {
			p.fatal("unknown option %q", key)
		}

		opt.Set = true

		if opt.ValueName == "" {
			args = args[1:]
		} else {
			if len(args) < 2 {
				p.fatal("missing value for option %q", key)
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
		p.fatal("missing command")
	}

	cmd := p.command
	fullName := []string{}

	for len(args) > 0 {
		name := args[0]
		fullName = append(fullName, name)

		cmd2 := cmd.subcommands[name]
		if cmd2 == nil {
			if cmd != nil {
				break
			}

			break
		}

		cmd = cmd2
		args = args[1:]
	}

	if cmd == nil || cmd == p.command {
		p.fatal("unknown command %q", strings.Join(fullName, " "))
	}

	if len(cmd.subcommands) > 0 {
		p.fatal("missing subcommand(s)")
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
			p.fatal("missing argument(s)")
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
				p.fatal("too many arguments")
			}
		}
	} else {
		if len(args) > 0 {
			p.fatal("unexpected arguments")
		}
	}

	return args
}
