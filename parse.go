package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	pupIn            io.ReadCloser = os.Stdin
	pupMaxPrintLevel int           = -1
	pupPrintColor    bool          = false
	pupIndentString  string        = " "
	pupDisplayer     Displayer     = TreeDisplayer{}
)

func PrintHelp(w io.Writer, exitCode int) {
	helpString := `Usage
    pup [flags] [selectors] [optional display function]
Version
    %s
Flags
    -c --color         print result with color
    -f --file          file to read from
    -h --help          display this help
    -i --indent        number of spaces to use for indent or character
    -n --number        print number of elements selected
    -l --limit         restrict number of levels printed
    --version          display version
`
	fmt.Fprintf(w, helpString, VERSION)
	os.Exit(exitCode)
}

func ParseArgs() []string {
	cmds := ProcessFlags(os.Args[1:])
	return ParseCommands(strings.Join(cmds, " "))
}

// Process command arguments and return all non-flags.
func ProcessFlags(cmds []string) []string {
	var i int
	var err error
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Option '%s' requires an argument", cmds[i])
			os.Exit(2)
		}
	}()
	nonFlagCmds := make([]string, len(cmds))
	n := 0
	for i = 0; i < len(cmds); i++ {
		cmd := cmds[i]
		switch cmd {
		case "-c", "--color":
			pupPrintColor = true
		case "-f", "--file":
			filename := cmds[i+1]
			pupIn, err = os.Open(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(2)
			}
			i++
		case "-h", "--help":
			PrintHelp(os.Stdout, 0)
		case "-i", "--indent":
			indentLevel, err := strconv.Atoi(cmds[i+1])
			if err == nil {
				pupIndentString = strings.Repeat(" ", indentLevel)
			} else {
				pupIndentString = cmds[i+1]
			}
			i++
		case "-l", "--limit":
			pupMaxPrintLevel, err = strconv.Atoi(cmds[i+1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Argument for '%s' must be numeric\n", cmd)
				os.Exit(2)
			}
			i++
		case "--version":
			fmt.Println(VERSION)
			os.Exit(0)
		default:
			if cmd[0] == '-' {
				fmt.Fprintf(os.Stderr, "Unrecognized flag '%s'", cmd)
				os.Exit(2)
			}
			nonFlagCmds[n] = cmds[i]
			n++
		}
	}
	return nonFlagCmds[:n]
}

// Split a string with awareness for quoted text
func ParseCommands(cmdString string) []string {
	cmds := []string{}
	last, next, max := 0, 0, len(cmdString)
	for {
		// if we're at the end of the string, return
		if next == max {
			if next > last {
				cmds = append(cmds, cmdString[last:next])
			}
			return cmds
		}
		// evalute a rune
		c := cmdString[next]
		switch c {
		case ' ':
			if next > last {
				cmds = append(cmds, cmdString[last:next])
			}
			last = next + 1
		case '\'', '"':
			// for quotes, consume runes until the quote has ended
			quoteChar := c
			for {
				next++
				if next == max {
					fmt.Fprintf(os.Stderr, "Unmatched open quote (%c)\n", quoteChar)
					os.Exit(2)
				}
				if cmdString[next] == quoteChar {
					break
				}
			}
		}
		next++
	}
}
