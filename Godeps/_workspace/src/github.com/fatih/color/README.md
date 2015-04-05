# Color [![GoDoc](https://godoc.org/github.com/fatih/color?status.svg)](http://godoc.org/github.com/fatih/color) [![Build Status](https://travis-ci.org/fatih/color.svg)](https://travis-ci.org/fatih/color)



Color let you use colorized outputs in terms of [ANSI Escape
Codes](http://en.wikipedia.org/wiki/ANSI_escape_code#Colors). The API can be
used in several way, pick one that suits you.



![Color](http://i.imgur.com/c1JI0lA.png)


## Install

```bash
go get github.com/fatih/color
```

## Examples

### Standard colors

```go
// Print with default helper functions
color.Cyan("Prints text in cyan.")

// a newline will be appended automatically
color.Blue("Prints %s in blue.", "text")

// These are using by default foreground colors.
color.Red("We have red")
color.Magenta("And many others ..")

```

### Mix and reuse colors

```go
// Create a new color object
c := color.New(color.FgCyan).Add(color.Underline)
c.Println("Prints cyan text with an underline.")

// Or just add them to New()
d := color.New(color.FgCyan, color.Bold)
d.Printf("This prints bold cyan %s\n", "too!.")

// Mix up foreground and background colors, create new mixes!
red := color.New(color.FgRed)

boldRed := red.Add(color.Bold)
boldRed.Println("This will print text in bold red.")

whiteBackground := red.Add(color.BgWhite)
whiteBackground.Println("Red text with White background.")
```

### Custom print functions (PrintFunc)

```go
// Create a custom print function for convenient
red := color.New(color.FgRed).PrintfFunc()
red("warning")
red("error: %s", err)

// Mix up multiple attributes
notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
notice("don't forget this...")
```

### Insert into noncolor strings (SprintFunc)

```go
// Create SprintXxx functions to mix strings with other non-colorized strings:
yellow := color.New(color.FgYellow).SprintFunc()
red := color.New(color.FgRed).SprintFunc()
fmt.Printf("this is a %s and this is %s.\n", yellow("warning"), red("error"))

info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
fmt.Printf("this %s rocks!\n", info("package"))

// use helper functions
fmt.Printf("this", color.RedString("warning"), "should be not neglected.")
fmt.Printf(color.GreenString("Info:"), "an important message." )
```

### Plug into existing code

```go
// Use handy standard colors.
color.Set(color.FgYellow)

fmt.Println("Existing text will be now in Yellow")
fmt.Printf("This one %s\n", "too")

color.Unset() // don't forget to unset

// You can mix up parameters
color.Set(color.FgMagenta, color.Bold)
defer color.Unset() // use it in your function

fmt.Println("All text will be now bold magenta.")
```

## Todo

* Save/Return previous values
* Add Windows support
* Evaluate fmt.Formatter interface


## Credits

 * [Fatih Arslan](https://github.com/fatih)

## License

The MIT License (MIT) - see LICENSE.md for more details


