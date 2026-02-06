package program

import (
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"go.n16f.net/uuid"
)

var commandNameRE = regexp.MustCompile("\\s+")

type Command struct {
	Name        string
	FullName    string
	Description string
	Main        Main

	program *Program

	subcommands map[string]*Command
	options     map[string]*Option
	arguments   []*Argument
}

func (p *Program) newCommandGroup(name, fullName string) *Command {
	return &Command{
		Name:     name,
		FullName: fullName,

		program: p,

		subcommands: make(map[string]*Command),
	}
}

func (c *Command) Label() string {
	if len(c.subcommands) == 0 {
		return c.Name
	}

	return c.Name + " <subcommand>..."
}

type Option struct {
	ShortName    string
	LongName     string
	ValueName    string
	DefaultValue string
	Description  string

	Set   bool
	Value string
}

type Argument struct {
	Name        string
	Description string
	Optional    bool
	Trailing    bool

	Set            bool
	Value          string
	TrailingValues []string
}

func (p *Program) AddCommand(fullName, description string, main Main) *Command {
	if p.Main != nil {
		panic("cannot have a main function with commands")
	}

	names := splitCommandName(fullName)
	if len(names) == 0 {
		panic("empty command name")
	}

	if p.command == nil {
		p.command = p.newCommandGroup("", "")
	}

	group := p.command

	for i := range len(names) - 1 {
		name := names[i]

		if group.subcommands == nil {
			group.subcommands = make(map[string]*Command)
		}

		group2 := group.subcommands[name]
		if group2 == nil {
			group2 = p.newCommandGroup(name, strings.Join(names[:i+1], " "))
			group.subcommands[name] = group2
		}

		group = group2
	}

	name := names[len(names)-1]

	cmd := Command{
		Name:        name,
		FullName:    strings.Join(names, " "),
		Description: description,
		Main:        main,

		program: p,

		options: make(map[string]*Option),
	}

	if cmd := group.subcommands[name]; cmd != nil {
		if len(cmd.subcommands) == 0 {
			Panic("duplicate command %q", cmd.FullName)
		} else {
			Panic("command %q has subcommands", cmd.FullName)
		}
	}

	if group.Main != nil {
		Panic("command %q cannot be used as a group", group.FullName)
	}

	group.subcommands[name] = &cmd

	return &cmd
}

func (p *Program) AddOption(shortName, longName, valueName, defaultValue, description string) {
	option := &Option{
		ShortName:    shortName,
		LongName:     longName,
		ValueName:    valueName,
		DefaultValue: defaultValue,
		Description:  description,
	}

	p.addOption(nil, option)
}

func (p *Program) AddFlag(shortName, longName, description string) {
	p.AddOption(shortName, longName, "", "", description)
}

func (c *Command) AddOption(shortName, longName, valueName, defaultValue, description string) {
	option := &Option{
		ShortName:    shortName,
		LongName:     longName,
		ValueName:    valueName,
		DefaultValue: defaultValue,
		Description:  description,
	}

	c.program.addOption(c, option)
}

func (c *Command) AddFlag(shortName, longName, description string) {
	c.AddOption(shortName, longName, "", "", description)
}

func (p *Program) addOption(c *Command, option *Option) {
	var m map[string]*Option

	if option.ShortName == "" && option.LongName == "" {
		panic("option has no short or long name")
	}

	if c == nil {
		m = p.options
	} else {
		m = c.options
	}

	if option.ShortName != "" {
		if _, found := m[option.ShortName]; found {
			Panic("duplicate option name %q", option.ShortName)
		}

		if c != nil {
			if _, found := c.program.options[option.ShortName]; found {
				Panic("duplicate option name %q", option.ShortName)
			}
		}

		m[option.ShortName] = option
	}

	if option.LongName != "" {
		if _, found := m[option.LongName]; found {
			Panic("duplicate option name %q", option.LongName)
		}

		if c != nil {
			if _, found := c.program.options[option.LongName]; found {
				Panic("duplicate option name %q", option.LongName)
			}
		}

		m[option.LongName] = option
	}
}

func (p *Program) AddArgument(name, description string) {
	checkForArgument(p.arguments)

	arg := &Argument{
		Name:        name,
		Description: description,
	}

	p.arguments = append(p.arguments, arg)
}

func (p *Program) AddOptionalArgument(name, description string) {
	checkForOptionalArgument(p.arguments)

	arg := &Argument{
		Name:        name,
		Description: description,
		Optional:    true,
	}

	p.arguments = append(p.arguments, arg)
}

func (p *Program) AddTrailingArgument(name, description string) {
	checkForTrailingArgument(p.arguments)

	arg := &Argument{
		Name:        name,
		Description: description,
		Trailing:    true,
	}

	p.arguments = append(p.arguments, arg)
}

func (c *Command) AddArgument(name, description string) {
	checkForArgument(c.arguments)

	arg := &Argument{
		Name:        name,
		Description: description,
	}

	c.arguments = append(c.arguments, arg)
}

func (c *Command) AddOptionalArgument(name, description string) {
	checkForOptionalArgument(c.arguments)

	arg := &Argument{
		Name:        name,
		Description: description,
		Optional:    true,
	}

	c.arguments = append(c.arguments, arg)
}

func (c *Command) AddTrailingArgument(name, description string) {
	checkForTrailingArgument(c.arguments)

	arg := &Argument{
		Name:        name,
		Description: description,
		Trailing:    true,
	}

	c.arguments = append(c.arguments, arg)
}

func checkForArgument(args []*Argument) {
	if len(args) == 0 {
		return
	}

	lastArg := args[len(args)-1]

	if lastArg.Optional {
		panic("cannot add non-optional argument after optional argument")
	}

	if lastArg.Trailing {
		panic("cannot add argument after trailing argument")
	}
}

func checkForOptionalArgument(args []*Argument) {
	if len(args) == 0 {
		return
	}

	lastArg := args[len(args)-1]

	if lastArg.Trailing {
		panic("cannot add argument after trailing argument")
	}
}

func checkForTrailingArgument(args []*Argument) {
	if len(args) == 0 {
		return
	}

	lastArg := args[len(args)-1]

	if lastArg.Trailing {
		panic("cannot add multiple trailing arguments")
	}
}

func (p *Program) CommandName() string {
	if p.command == nil {
		Panic("no command defined")
	}

	return p.selectedCommand.Name
}

func (p *Program) CommandFullName() string {
	if p.command == nil {
		Panic("no command defined")
	}

	return p.selectedCommand.FullName
}

func (p *Program) IsOptionSet(name string) bool {
	return p.mustOption(name).Set
}

func (p *Program) OptionValue(name string) string {
	opt := p.mustOption(name)
	if !opt.Set {
		return opt.DefaultValue
	}

	return opt.Value
}

func (p *Program) BooleanOptionValue(name string) bool {
	return p.booleanValue("option", name, p.OptionValue(name))
}

func (p *Program) UUIDOptionValue(name string) uuid.UUID {
	return p.uuidValue("option", name, p.OptionValue(name))
}

func (p *Program) mustOption(name string) *Option {
	if cmd := p.selectedCommand; cmd != nil {
		option, found := cmd.options[name]
		if found {
			return option
		}
	}

	option, found := p.options[name]
	if !found {
		Panic("unknown option %q", name)
	}

	return option
}

func (p *Program) ArgumentValue(name string) string {
	return p.mustArgument(name).Value
}

func (p *Program) BooleanArgumentValue(name string) bool {
	return p.booleanValue("argument", name, p.ArgumentValue(name))
}

func (p *Program) UUIDArgumentValue(name string) uuid.UUID {
	return p.uuidValue("argument", name, p.ArgumentValue(name))
}

func (p *Program) booleanValue(typeName, name, value string) bool {
	switch strings.ToLower(value) {
	case "true":
		return true
	case "false":
		return false
	}

	p.Fatal("invalid value %q for %s %q: must be either %q or %q",
		value, typeName, name, "true", "false")
	return false // make the Go compiler happy
}

func (p *Program) uuidValue(typeName, name, value string) uuid.UUID {
	var id uuid.UUID
	if err := id.Parse(value); err != nil {
		p.Fatal("invalid value %q for %s %q: must be a valid UUID",
			value, typeName, name)
	}

	return id
}

func (p *Program) OptionalArgumentValue(name string) *string {
	arg := p.mustArgument(name)
	if !arg.Set {
		return nil
	}

	return &arg.Value
}

func (p *Program) TrailingArgumentValues(name string) []string {
	return p.mustArgument(name).TrailingValues
}

func (p *Program) mustArgument(name string) *Argument {
	var arguments []*Argument

	if cmd := p.selectedCommand; cmd == nil {
		arguments = p.arguments
	} else {
		arguments = cmd.arguments
	}

	for _, argument := range arguments {
		if name == argument.Name {
			return argument
		}
	}

	Panic("unknown argument %q", name)
	return nil // make the compiler happy
}

func (p *Program) findCommand(names []string) *Command {
	cmd := p.command

	for _, name := range names {
		cmd = cmd.subcommands[name]
		if cmd == nil {
			return nil
		}
	}

	if cmd == p.command {
		return nil
	}

	return cmd
}

func (p *Program) ParseCommandLine() {
	if p.command != nil {
		p.addDefaultCommands()
	}

	p.parse()

	if p.IsOptionSet("help") {
		cmdHelp(p)
		os.Exit(0)
	}

	p.Quiet = p.IsOptionSet("quiet")

	if p.IsOptionSet("debug") {
		s := p.OptionValue("debug")
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil || i < 0 || i > math.MaxInt32 {
			p.Fatal("invalid debug level %q", s)
		}

		p.DebugLevel = int(i)
	}
}

func (p *Program) addDefaultOptions() {
	p.AddFlag("h", "help", "print help and exit")
	p.AddFlag("q", "quiet", "do not print status and information messages")
	p.AddOption("", "debug", "level", "0", "print debug messages")
}

func (p *Program) addDefaultCommands() {
	c := p.AddCommand("help", "print help and exit", cmdHelp)
	c.AddTrailingArgument("command", "the name of the command")
}

func cmdHelp(p *Program) {
	cmd := p.command

	if p.selectedCommand != nil {
		if p.selectedCommand.FullName == "help" && !p.IsOptionSet("help") {
			names := p.TrailingArgumentValues("command")
			if len(names) > 0 {
				cmd = p.findCommand(names)
				if cmd == nil {
					p.Fatal("unknown command %q", strings.Join(names, " "))
				}
			}
		} else {
			cmd = p.selectedCommand
		}
	}

	p.PrintUsage(cmd)
}

func splitCommandName(s string) []string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return []string{}
	}

	parts := commandNameRE.Split(s, -1)

	var names []string
	for _, part := range parts {
		if name := strings.TrimSpace(part); name != "" {
			names = append(names, name)
		}
	}

	return names
}
