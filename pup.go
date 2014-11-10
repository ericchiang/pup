package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

//      _=,_
//   o_/6 /#\
//   \__ |##/
//    ='|--\
//      /   #'-.
//      \#|_   _'-. /
//       |/ \_( # |"
//      C/ ,--___/

var VERSION string = "0.3.4"

func main() {
	// process flags and arguments
	cmds := ParseArgs()

	// Determine the charset of the input
	cr, err := charset.NewReader(pupIn, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}

	// Parse the input and get the root node
	root, err := html.Parse(cr)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}

	// Parse the selectors
	selectorFuncs := []SelectorFunc{}
	funcGenerator := Select
	var cmd string
	for len(cmds) > 0 {
		cmd, cmds = cmds[0], cmds[1:]
		if len(cmds) == 0 {
			if err := ParseDisplayer(cmd); err == nil {
				continue
			}
		}
		switch cmd {
		case "*":
			continue
		case "+":
			funcGenerator = SelectFromChildren
		case ">":
			funcGenerator = SelectNextSibling
		default:
			selector, err := ParseSelector(cmd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Selector parsing error: %s\n", err.Error())
				os.Exit(2)
			}
			selectorFuncs = append(selectorFuncs, funcGenerator(selector))
			funcGenerator = Select
		}
	}

	currNodes := []*html.Node{root}
	for _, selectorFunc := range selectorFuncs {
		currNodes = selectorFunc(currNodes)
	}
	pupDisplayer.Display(currNodes)
}
