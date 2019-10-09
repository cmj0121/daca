package args

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	/* Library version */
	MAJOR = 1
	MINOR = 0
	MACRO = 0
)

type option interface {
	/* option setter*/
	Shortcut(in byte) option                                 /* set the shortcut of the option        */
	Default(in interface{}) option                           /* set the default value                 */
	Choice(in []string) option                               /* set the possible values               */
	Required(in bool) option                                 /* set option required or NOT            */
	Validator(fn func(interface{}) error) option             /* validator of the user's input         */
	Callback(fn func(p *argparser, opt option) error) option /* callback function when option matched */
	Help(in string) option                                   /* set the help message                  */
	/* option getter */
	Get() interface{}                /* get the option value          */
	Set(in []string, idx *int) error /* set the value to option       */
	String(format string) string     /* the help message              */
	Hint() string                    /* the hint (short) of flag      */
	Check() error                    /* check everything is OK        */
}

type argparser struct {
	Name    string
	Tooltip string

	version   string            /* internal property                       */
	options   map[string]option /* options hash-map  via the name          */
	shortcut  map[byte]option   /* shortcut of the options                 */
	parameter []option          /* extra parameters                        */
	opt_list  []option          /* the inner option list for help message  */
}

func Parser(tooltip string) (out *argparser) {
	out = SubParser("", tooltip)
	return
}

func SubParser(name, tooltip string) (out *argparser) {
	out = &argparser{
		Name:      name,
		Tooltip:   tooltip,
		options:   make(map[string]option, 0),
		shortcut:  make(map[byte]option, 0),
		parameter: make([]option, 0),
	}

	/* set the default option */
	help := out.Flip("help").Shortcut('h').Help("Show this message")
	help.Callback(
		func(p *argparser, opt option) (err error) {
			fmt.Println(p)
			os.Exit(1)
			return
		},
	)

	version := out.Flip("version").Shortcut('V').Help("Show version")
	version.Callback(
		func(p *argparser, opt option) (err error) {
			fmt.Println(fmt.Sprintf("%s (%s)", out.Name, out.version))
			os.Exit(0)
			return
		},
	)

	return
}

func (p *argparser) Version(ver string) (out *argparser) {
	out = p
	out.version = ver
	return
}

func (p *argparser) Get(key string) (out interface{}) {
	if v, ok := p.options[key]; ok {
		out = v.Get()
	}

	return
}

func (p *argparser) Shortcut(shortcut byte, in option) (out option) {
	if _, ok := p.shortcut[shortcut]; ok {
		panic(fmt.Sprintf("duplicate shortcut '%c'", shortcut))
	}

	p.shortcut[shortcut] = in
	out = in
	return
}

func (p *argparser) Flip(name string) (out option) {
	out = &Flip{
		name:   name,
		parser: p,
	}

	if _, ok := p.options[name]; ok {
		panic(fmt.Sprintf("duplicate option name '%s'", name))
	}

	p.options[name] = out
	p.opt_list = append(p.opt_list, out)
	return
}

func (p *argparser) Flag(name string) (out option) {
	out = &Flag{
		name:   name,
		parser: p,
	}

	if _, ok := p.options[name]; ok {
		panic(fmt.Sprintf("duplicate option name '%s'", name))
	}

	p.options[name] = out
	p.opt_list = append(p.opt_list, out)
	return
}

func (p *argparser) Parameter(name string) (out option) {
	out = &Parameter{
		name:   name,
		parser: p,
	}

	if _, ok := p.options[name]; ok {
		panic(fmt.Sprintf("duplicate option name '%s'", name))
	}

	p.options[name] = out
	p.parameter = append(p.parameter, out)
	p.opt_list = append(p.opt_list, out)
	return
}

func (p *argparser) Sub(name string, subs ...*argparser) (out option) {
	choices := make([]string, 0)
	subparser := make(map[string]*argparser, 0)

	for _, sub := range subs {
		choices = append(choices, sub.Name)

		if _, ok := subparser[sub.Name]; ok {
			panic(fmt.Sprintf("duplicate sub-command '%s'", sub.Name))
		}
		subparser[sub.Name] = sub
	}
	out = &Sub{
		name:      name,
		parser:    p,
		choices:   choices,
		subparser: subparser,
	}

	if _, ok := p.options[name]; ok {
		panic(fmt.Sprintf("duplicate option name '%s'", name))
	}

	p.options[name] = out
	p.parameter = append(p.parameter, out)
	p.opt_list = append(p.opt_list, out)
	return
}

func (p *argparser) String() (out string) {
	var lines []string

	if p.version != "" {
		lines = append(lines, fmt.Sprintf("%s (%s) : %s", p.Name, p.version, p.Tooltip))
	} else {
		lines = append(lines, fmt.Sprintf("%s : %s", p.Name, p.Tooltip))
	}

	hints := make([]string, 0)
	for _, opt := range p.opt_list {
		hints = append(hints, opt.Hint())
	}

	lines = append(lines, fmt.Sprintf("usage: %s %s", p.Name, strings.Join(hints, " ")))
	if p.options != nil {
		lines = append(lines, "")
		for _, opt := range p.opt_list {
			lines = append(lines, opt.String("  %-14s %s"))
		}
	}

	out = strings.Join(lines, "\n")
	return
}

func (p *argparser) Exit(code int, err error) {
	fmt.Println(fmt.Sprintf("error: %s\n", err))
	fmt.Println(p)
	os.Exit(code)
}

func (p *argparser) Parse(args ...string) (err error) {
	p.Name = args[0]

	for idx := 1; idx < len(args); idx++ {
		argv := args[idx]

		switch {
		case len(argv) > 2 && "--" == argv[:2]:
			key := strings.Split(argv[2:], "=")[0]

			if opt, ok := p.options[key]; ok {
				if err = opt.Set(args, &idx); err != nil {
					return
				}
				continue
			}
			err = errors.New(fmt.Sprintf("Not support options '%s'", key))
			return
		case len(argv) >= 2 && '-' == argv[0]:
			key := argv[1]

			if opt, ok := p.shortcut[key]; ok {
				if err = opt.Set(args, &idx); err != nil {
					return
				}
				continue
			}
			err = errors.New(fmt.Sprintf("Not support options '%c'", key))
			return
		default:
			for _, opt := range p.parameter {
				if opt.Get() == nil {
					if err = opt.Set(args, &idx); err != nil {
						return
					}
					break
				}
			}
		}
	}

	/* check the required */
	for _, opt := range p.options {
		if err = opt.Check(); err != nil {
			return
		}
	}
	return
}

/* general Choice validator */
func ChoiceValidator(in string, choices []string) (err error) {
	if choices == nil {
		return
	}

	for _, choice := range choices {
		if in == choice {
			return
		}
	}

	err = errors.New(fmt.Sprintf("value '%s' not in Choice %v", in, choices))
	return
}
