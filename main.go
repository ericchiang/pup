package main

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"fmt"
	"github.com/ericchiang/pup/selector"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Flags
	inputStream   io.ReadCloser = os.Stdin
	sep           string        = " "
	maxPrintLevel int           = -1
	printNumber   bool          = false

	// Helpers
	whitespaceRegexp *regexp.Regexp = regexp.MustCompile(`^\s*$`)
	preWhitespace    *regexp.Regexp = regexp.MustCompile(`^\s+`)
	postWhitespace   *regexp.Regexp = regexp.MustCompile(`\s+$`)
)

func printIndent(level int) {
	for ; level > 0; level-- {
		fmt.Print(sep)
	}
}

// Is this node a tag with no end tag such as <meta> or <br>?
// http://www.w3.org/TR/html-markup/syntax.html#syntax-elements
func isVoidElement(n *html.Node) bool {
	switch n.DataAtom {
	case atom.Area, atom.Base, atom.Br, atom.Col, atom.Command, atom.Embed,
		atom.Hr, atom.Img, atom.Input, atom.Keygen, atom.Link,
		atom.Meta, atom.Param, atom.Source, atom.Track, atom.Wbr:
		return true
	}
	return false

}

func printChildren(n *html.Node, level int) {
	if maxPrintLevel > -1 {
		if level >= maxPrintLevel {
			printIndent(level)
			fmt.Println("...")
			return
		}
	}
	child := n.FirstChild
	for child != nil {
		PrintNode(child, level)
		child = child.NextSibling
	}
}

func PrintNode(n *html.Node, level int) {
	switch n.Type {
	case html.TextNode:
		s := n.Data
		if !whitespaceRegexp.MatchString(s) {
			s = preWhitespace.ReplaceAllString(s, "")
			s = postWhitespace.ReplaceAllString(s, "")
			printIndent(level)
			fmt.Println(s)
		}
	case html.ElementNode:
		printIndent(level)
		fmt.Printf("<%s", n.Data)
		for _, a := range n.Attr {
			fmt.Printf(` %s="%s"`, a.Key, a.Val)
		}
		fmt.Print(">\n")
		if !isVoidElement(n) {
			printChildren(n, level+1)
			printIndent(level)
			fmt.Printf("</%s>\n", n.Data)
		}
	case html.CommentNode, html.DoctypeNode, html.DocumentNode:
		printChildren(n, level)
	}
}

func Fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func printHelp() {
	Fatal(`Usage:

    pup [list of css selectors]

Flags:

    -f --file          file to read from
    -h --help          display this help
    -i --indent        number of spaces to use for indent or character
    -n --number        print number of elements selected
    -l --level         restrict number of levels printed 
`)
}

func processFlags(cmds []string) []string {
	var i int
	var err error
	defer func() {
		if r := recover(); r != nil {
			Fatal("Option '%s' requires an argument", cmds[i])
		}
	}()
	nonFlagCmds := make([]string, len(cmds))
	n := 0
	for i = 0; i < len(cmds); i++ {
		cmd := cmds[i]
		switch cmd {
		case "-f", "--file":
			filename := cmds[i+1]
			inputStream, err = os.Open(filename)
			if err != nil {
				Fatal(err.Error())
			}
			i++
		case "-h", "--help":
			printHelp()
			os.Exit(1)
		case "-i", "--indent":
			indentLevel, err := strconv.Atoi(cmds[i+1])
			if err == nil {
				sep = strings.Repeat(" ", indentLevel)
			} else {
				sep = cmds[i+1]
			}
			i++
		case "-n", "--number":
			printNumber = true
		case "-l", "--level":
			maxPrintLevel, err = strconv.Atoi(cmds[i+1])
			if err != nil {
				Fatal("Argument for '%s' must be numeric",
					cmds)
			}
			i++
		default:
			if cmd[0] == '-' {
				Fatal("Unrecognized flag '%s'", cmd)
			}
			nonFlagCmds[n] = cmds[i]
			n++
		}
	}
	return nonFlagCmds[:n]
}

func main() {
	cmds := processFlags(os.Args[1:])
	root, err := html.Parse(inputStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
	inputStream.Close()
	if len(cmds) == 0 {
		PrintNode(root, 0)
		os.Exit(0)
	}
	selectors := make([]selector.Selector, len(cmds))
	for i, cmd := range cmds {
		selectors[i], err = selector.ParseSelector(cmd)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(2)
		}
	}
	currNodes := []*html.Node{root}
	var selected []*html.Node
	for _, selector := range selectors {
		selected = []*html.Node{}
		for _, node := range currNodes {
			selected = append(selected, selector.FindAllChildren(node)...)
		}
		currNodes = selected
	}
	if printNumber {
		fmt.Println(len(currNodes))
	} else {
		for _, s := range currNodes {
			PrintNode(s, 0)
		}
	}
}
