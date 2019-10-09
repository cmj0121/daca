# args #
The args is the simple argument parser library for go. The simple code of args is as following

```go
package main

import (
	"github.com/cmj0121/args"
	"fmt"
	"os"
)

func main() {
	p := args.Parser("simple argument parser").Version("alpha")

	p.Flip("enable").Shortcut('E').Required(true).Help("The flip option")
	p.Flag("name").Shortcut('n').Default("cmj").Help("The flag option")

	if err := p.Parse(os.Args...); err != nil {
		p.Exit(1, err)
	}

	fmt.Println("final options")
	fmt.Println(fmt.Sprintf("- enable : %v", p.Get("enable")))
	fmt.Println(fmt.Sprintf("- name   : %v", p.Get("name")))
}
```

## options ##
In the args, there are three kinds of the options: Flip, Flag and Parameter. The flip option is used to
store the true/false value and it will flip if there are multiple input in the command line. For example,
there is a flip option, called *-E* which default value is the **true**. If executed with one *-E* in command
line, the option is stored the *false* value but it will become **true** again if type two *-E* in command
line.

The flag option is used to store extra value from command line. It support three kinds of format to provide
the value: 1) the following remainder string when using the shortcut, 2) the following string of the option
and 3) the connected string via the equal sign. The following is the sample command line of the options *-n, --name*:

```bash
> example -ncmj      # The first case for the shortcut
> example --name cmj # The following value of the options
> example --name=cmj # The connected value via the equal sign
```

## sub-command ##
In some case the programmer need build a command-line tool with sub-command: the second argument parser for
the inner command. For example, the main tool has a option with two possible value, setter and getter, and
each value has related flow for other options. In the args it provide the *Sub* option for the sub-command
and using the *SubParser*  to build the nested parser for this sub-command.


```go
package main

import (
	"args"
	"fmt"
	"os"
)

func main() {
	p := args.Parser("simple sub-command parser").Version("beta")

	sub_setter := args.SubParser("setter", "sub-command, setter")
	sub_getter := args.SubParser("getter", "sub-command, getter")
	p.Sub("action", sub_setter, sub_getter).Help("Sub-command")

	if err := p.Parse(os.Args...); err != nil {
		p.Exit(1, err)
	}
}
```
