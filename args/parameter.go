package args

import (
	"errors"
	"fmt"
	"strings"
)

type Parameter struct {
	name   string
	parser *argparser

	value         interface{}
	default_value interface{}
	choices       []string
	required      bool
	validator     func(interface{}) error
	callback      func(p *argparser, opt option) error
	help          string
}

/* getter */
func (opt *Parameter) Get() (out interface{}) {
	out = opt.value
	if out == nil {
		out = opt.default_value
	}
	return
}

func (opt *Parameter) Set(in []string, idx *int) (err error) {
	value := in[*idx]
	*idx++

	if opt.validator != nil {
		if err = opt.validator(value); err != nil {
			return
		}
	}

	/* choices validator */
	if err = ChoiceValidator(value, opt.choices); err != nil {
		return
	}

	if opt.callback != nil {
		err = opt.callback(opt.parser, opt)
	}

	opt.value = value
	return
}

func (opt *Parameter) String(format string) (out string) {
	out = fmt.Sprintf(format, opt.name, opt.help)

	if opt.default_value != nil {
		out = fmt.Sprintf("%s (default: %v)", out, opt.default_value)
	}

	if opt.choices != nil {
		out = fmt.Sprintf("%s %v", out, opt.choices)
	}

	return
}

func (opt *Parameter) Hint() (out string) {
	if opt.required == true {
		out = fmt.Sprintf("%s", strings.ToUpper(opt.name))
	} else {
		out = fmt.Sprintf("[%s]", strings.ToUpper(opt.name))
	}

	return
}

func (opt *Parameter) Check() (err error) {
	if opt.required && opt.value == nil {
		err = errors.New(fmt.Sprintf("options '%s' is required", opt.name))
	}
	return
}

/* setter */
func (opt *Parameter) Shortcut(in byte) (out option) {
	panic("Parameter does NOT support Shortcut")
	return
}

func (opt *Parameter) Default(in interface{}) (out option) {
	opt.default_value = in.(string)
	out = opt
	return
}

func (opt *Parameter) Choice(in []string) (out option) {
	opt.choices = in
	out = opt
	return
}

func (opt *Parameter) Required(in bool) (out option) {
	opt.required = in
	out = opt
	return
}

func (opt *Parameter) Validator(in func(interface{}) error) (out option) {
	opt.validator = in
	out = opt
	return
}

func (opt *Parameter) Callback(in func(p *argparser, opt option) error) (out option) {
	opt.callback = in
	out = opt
	return
}

func (opt *Parameter) Help(in string) (out option) {
	opt.help = in
	out = opt
	return
}
