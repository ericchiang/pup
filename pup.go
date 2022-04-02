package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

//      _=,_
//   o_/6 /#\
//   \__ |##/
//    ='|--\
//      /   #'-.
//      \#|_   _'-. /
//       |/ \_( # |"
//      C/ ,--___/

var VERSION string = "0.4.0"

func main() {
	// process flags and arguments
	cmds, err := ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	// Parse the input and get the root node
	root, err := ParseHTML(pupIn, pupCharset)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}
	pupIn.Close()

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
		case "*": // select all
			continue
		case ">":
			funcGenerator = SelectFromChildren
		case "+":
			funcGenerator = SelectNextSibling
		case ",": // nil will signify a comma
			selectorFuncs = append(selectorFuncs, nil)
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

	selectedNodes := []*html.Node{}
	currNodes := []*html.Node{root}
	for _, selectorFunc := range selectorFuncs {
		if selectorFunc == nil { // hit a comma
			selectedNodes = append(selectedNodes, currNodes...)
			currNodes = []*html.Node{root}
		} else {
			currNodes = selectorFunc(currNodes)
		}
	}
	selectedNodes = append(selectedNodes, currNodes...)
	if pupInvertSelect {
		RemoveInverseMatches(root, selectedNodes)
		pupDisplayer.Display([]*html.Node{root})
	} else {
		pupDisplayer.Display(selectedNodes)
	}
}

func RemoveInverseMatches(root *html.Node, selectedNodes []*html.Node) {
	var remove []*html.Node
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		if HtmlNodeInList(node, selectedNodes) {
			remove = append(remove, node)
		} else {
			RemoveInverseMatches(node, selectedNodes)
		}
	}
	for _, rm := range remove {
		root.RemoveChild(rm)
	}
}

func HtmlNodeInList(a *html.Node, list []*html.Node) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
