package program

import (
	"math"
	"os"
	"strconv"
)

type Command struct {
	Name        string
	Description string
	Main        Main

	program *Program

	options   map[string]*Option
	arguments []*Argument
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

func (p *Program) AddCommand(name, description string, main Main) *Command {
	if p.Main != nil {
		panic("cannot have a main function with commands")
	}

	c := &Command{
		Name:        name,
		Description: description,
		Main:        main,

		program: p,

		options: make(map[string]*Option),
	}

	p.commands[name] = c

	return c
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
	if len(p.commands) == 0 {
		Panic("no command defined")
	}

	return p.command.Name
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

func (p *Program) mustOption(name string) *Option {
	if p.command != nil {
		option, found := p.command.options[name]
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

	if p.command == nil {
		arguments = p.arguments
	} else {
		arguments = p.command.arguments
	}

	for _, argument := range arguments {
		if name == argument.Name {
			return argument
		}
	}

	Panic("unknown argument %q", name)
	return nil // make the compiler happy
}

func (p *Program) ParseCommandLine() {
	if len(p.commands) > 0 {
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
			p.fatal("invalid debug level %v", s)
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
	c.AddOptionalArgument("command", "the name of the command")
}

func cmdHelp(p *Program) {
	var commandName *string
	if p.command != nil {
		if p.command.Name == "help" {
			commandName = p.OptionalArgumentValue("command")
		} else {
			commandName = &p.command.Name
		}
	}

	if commandName == nil {
		p.PrintUsage(nil)
	} else {
		command, found := p.commands[*commandName]
		if !found {
			p.Error("unknown command %q", *commandName)
			os.Exit(1)
		}

		p.PrintUsage(command)
	}
}
