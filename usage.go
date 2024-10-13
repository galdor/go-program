package program

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"

	"golang.org/x/exp/maps"
)

func (p *Program) PrintUsage(cmd *Command) {
	var buf bytes.Buffer

	programName := os.Args[0]
	if cmd != nil && cmd.FullName != "" {
		programName += " " + cmd.FullName
	}

	var hasCommands bool
	var commands map[string]*Command
	if cmd != nil {
		hasCommands = len(cmd.subcommands) > 0
		commands = cmd.subcommands
	}

	var arguments []*Argument
	var description string

	if cmd == nil {
		arguments = p.arguments
		description = p.Description
	} else {
		arguments = cmd.arguments
		description = cmd.Description
	}

	hasArguments := len(arguments) > 0
	maxWidth := p.computeMaxWidth(cmd)

	fmt.Fprintf(&buf, "Usage: %s OPTIONS", programName)

	if hasCommands && !hasArguments {
		fmt.Fprintf(&buf, " <command>...")
	}

	if hasArguments {
		for _, arg := range arguments {
			if arg.Trailing {
				fmt.Fprintf(&buf, " [<%s>...]", arg.Name)
			} else if arg.Optional {
				fmt.Fprintf(&buf, " [<%s>]", arg.Name)
			} else {
				fmt.Fprintf(&buf, " <%s>", arg.Name)
			}
		}
	}

	fmt.Fprintf(&buf, "\n")

	if description != "" {
		fmt.Fprintf(&buf, "\n%s\n", sentence(description))
	}

	if hasCommands {
		p.usageCommands(&buf, commands, maxWidth)
	} else if hasArguments {
		p.usageArguments(&buf, arguments, maxWidth)
	}

	if len(p.options) > 0 {
		if cmd != nil && len(cmd.options) > 0 {
			p.usageOptions(&buf, "GLOBAL OPTIONS", p.options, maxWidth)
		} else {
			p.usageOptions(&buf, "OPTIONS", p.options, maxWidth)
		}
	}

	if cmd != nil && len(cmd.options) > 0 {
		p.usageOptions(&buf, "COMMAND OPTIONS", cmd.options, maxWidth)
	}

	io.Copy(os.Stderr, &buf)
}

func (p *Program) computeMaxWidth(cmd *Command) int {
	max := 0

	if cmd != nil {
		for _, subcmd := range cmd.subcommands {
			if label := subcmd.Label(); len(label) > max {
				max = len(label)
			}
		}
	}

	var args []*Argument
	if cmd == nil {
		args = p.arguments
	} else {
		args = cmd.arguments
	}

	for _, arg := range args {
		if len(arg.Name) > max {
			max = len(arg.Name)
		}
	}

	f := func(opt *Option) {
		length := 2 + 2 + 2 + len(opt.LongName)
		if opt.ValueName != "" {
			length += 2 + len(opt.ValueName) + 1
		}

		if length > max {
			max = length
		}
	}

	for _, opt := range p.options {
		f(opt)
	}

	if cmd != nil {
		for _, opt := range cmd.options {
			f(opt)
		}
	}

	return max
}

func (p *Program) usageCommands(buf *bytes.Buffer, commands map[string]*Command, maxWidth int) {
	fmt.Fprintf(buf, "\nCOMMANDS\n\n")

	names := maps.Keys(commands)
	slices.Sort(names)

	for _, name := range names {
		cmd := commands[name]
		fmt.Fprintf(buf, "%-*s  %s\n", maxWidth, cmd.Label(), cmd.Description)
	}
}

func (p *Program) usageArguments(buf *bytes.Buffer, args []*Argument, maxWidth int) {
	fmt.Fprintf(buf, "\nARGUMENTS\n\n")

	for _, arg := range args {
		fmt.Fprintf(buf, "%-*s  %s\n", maxWidth, arg.Name, arg.Description)
	}
}

func (p *Program) usageOptions(buf *bytes.Buffer, label string, options map[string]*Option, maxWidth int) {
	fmt.Fprintf(buf, "\n%s\n\n", label)

	strs := make(map[*Option]string)

	for _, opt := range options {
		if _, found := strs[opt]; found {
			continue
		}

		buf := bytes.NewBuffer([]byte{})

		if opt.ShortName == "" {
			fmt.Fprintf(buf, "  ")
		} else {
			fmt.Fprintf(buf, "-%s", opt.ShortName)
		}

		if opt.LongName != "" {
			if opt.ShortName == "" {
				buf.WriteString("  ")
			} else {
				buf.WriteString(", ")
			}

			fmt.Fprintf(buf, "--%s", opt.LongName)
		}

		if opt.ValueName != "" {
			fmt.Fprintf(buf, " <%s>", opt.ValueName)
		}

		str := buf.String()
		strs[opt] = str
	}

	var opts []*Option
	for opt, _ := range strs {
		opts = append(opts, opt)
	}

	sort.Slice(opts, func(i, j int) bool {
		return opts[i].sortKey() < opts[j].sortKey()
	})

	for _, opt := range opts {
		fmt.Fprintf(buf, "%-*s  %s", maxWidth, strs[opt], opt.Description)

		if opt.DefaultValue != "" {
			fmt.Fprintf(buf, " (default: %q)", opt.DefaultValue)
		}

		fmt.Fprintf(buf, "\n")
	}
}

func (opt *Option) sortKey() string {
	if opt.ShortName != "" {
		return opt.ShortName
	}

	if opt.LongName != "" {
		return opt.LongName
	}

	return ""
}
