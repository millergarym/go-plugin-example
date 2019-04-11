package main

import (
	"fmt"
	"os"
	"os/exec"
	"plugin"
)

type Greeter interface {
	Greet()
}

func build(l string) error {
	opath := "./" + l + "/" + l + ".so"
	spath := "./" + l + "/greeter.go"
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", opath, spath)
	return cmd.Run()
}

func main() {
	// determine module to load
	lang := []string{"eng"}
	if len(os.Args) >= 2 {
		lang = os.Args[1:]
	}

	plugins := make([]*plugin.Plugin, 0)
	for _, l := range lang {
		// load module
		// 1. open the so file to load the symbols
		plug, err := plugin.Open("./" + l + "/" + l + ".so")
		if err != nil {
			err = build(l)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plug, err = plugin.Open("./" + l + "/" + l + ".so")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		plugins = append(plugins, plug)
	}

	for _, plug := range plugins {
		// 2. look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		symGreeter, err := plug.Lookup("Greeter")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 3. Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		var greeter Greeter
		greeter, ok := symGreeter.(Greeter)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}

		// 4. use the module
		greeter.Greet()
	}
}
