package main

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/charset"
	"fmt"
	"github.com/ericchiang/pup/funcs"
	"github.com/ericchiang/pup/selector"
	"io"
	"os"
	"strconv"
	"strings"
)

const VERSION string = "0.3.0"

var (
	// Flags
	attributes    []string        = []string{}
	inputStream   io.ReadCloser   = os.Stdin
	indentString  string          = " "
	maxPrintLevel int             = -1
	printNumber   bool            = false
	printColor    bool            = false
	displayer     funcs.Displayer = nil
)

// Print to stderr and exit
func Fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

// Print help to stderr and quit
func PrintHelp() {
	helpString := `Usage

    pup [list of css selectors]

Version

    %s

Flags

    -c --color         print result with color
    -f --file          file to read from
    -h --help          display this help
    -i --indent        number of spaces to use for indent or character
    -n --number        print number of elements selected
    -l --limit         restrict number of levels printed
    --version          display version`
	Fatal(helpString, VERSION)
}

// Process command arguments and return all non-flags.
func ProcessFlags(cmds []string) []string {
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
		case "-a", "--attr":
			attributes = append(attributes, cmds[i+1])
			i++
		case "-c", "--color":
			printColor = true
		case "-f", "--file":
			filename := cmds[i+1]
			inputStream, err = os.Open(filename)
			if err != nil {
				Fatal(err.Error())
			}
			i++
		case "-h", "--help":
			PrintHelp()
			os.Exit(1)
		case "-i", "--indent":
			indentLevel, err := strconv.Atoi(cmds[i+1])
			if err == nil {
				indentString = strings.Repeat(" ", indentLevel)
			} else {
				indentString = cmds[i+1]
			}
			i++
		case "-n", "--number":
			printNumber = true
		case "-l", "--limit":
			maxPrintLevel, err = strconv.Atoi(cmds[i+1])
			if err != nil {
				Fatal("Argument for '%s' must be numeric",
					cmds)
			}
			i++
		case "--version":
			Fatal(VERSION)
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

// Split a string while ignoring strings. 
func QuoteSplit(input string) []string {
	last := 0
	split := []string{}
	inQuote := false
	quoteChar := ' '
	escapeNext := false
	for i, c := range input {
		if escapeNext {
			escapeNext = false
			continue
		}
		switch c {
		case ' ':
			if !inQuote {
				if last < i-1 {
					split = append(split, input[last:i])
				}
				last = i + 1
			}
		case '"', '\'':
			if inQuote {
				if c == quoteChar {
					inQuote = false
				}
			} else {
				inQuote = true
				quoteChar = c
			}
		case '\\':
			escapeNext = true
		}
	}
	if last < len(input) {
		split = append(split, input[last:len(input)])
	}
	return split
}

// pup
func main() {
	args := QuoteSplit(strings.Join(os.Args[1:], " "))
	cmds := ProcessFlags(args)
	cr, err := charset.NewReader(inputStream, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
	root, err := html.Parse(cr)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
	inputStream.Close()
	if len(cmds) == 0 {
		t := TreeDisplayer{indentString}
		t.Display([]*html.Node{root})
		os.Exit(0)
	}
	selectors := make([]selector.Selector, len(cmds))
	for i, cmd := range cmds {
		// if this is the last element, check for a function like
		// text{} or attr{}
		if i+1 == len(cmds) {
			d, err := funcs.NewDisplayFunc(cmd)
			if err == nil {
				displayer = d
				selectors = selectors[0 : len(cmds)-1]
				break
			} else {
				displayer = TreeDisplayer{indentString}
			}
		}
		selectors[i], err = selector.NewSelector(cmd)
		if err != nil {
			Fatal("Selector parse error: %s", err)
		}
	}
	currNodes := []*html.Node{root}
	for _, selector := range selectors {
		currNodes = selector.Select(currNodes)
	}
	if printNumber {
		fmt.Println(len(currNodes))
	} else {
		displayer.Display(currNodes)
	}
}
