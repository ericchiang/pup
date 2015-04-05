/*
Package color is an ANSI color package to output colorized or SGR defined
output to the standard output. The API can be used in several way, pick one
that suits you.

Use simple and default helper functions with predefined foreground colors:

    color.Cyan("Prints text in cyan.")

    // a newline will be appended automatically
    color.Blue("Prints %s in blue.", "text")

    // More default foreground colors..
    color.Red("We have red")
    color.Yellow("Yellow color too!")
    color.Magenta("And many others ..")

However there are times where custom color mixes are required. Below are some
examples to create custom color objects and use the print functions of each
separate color object.

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


You can create PrintXxx functions to simplify even more:

    // Create a custom print function for convenient
    red := color.New(color.FgRed).PrintfFunc()
    red("warning")
    red("error: %s", err)

    // Mix up multiple attributes
    notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
    notice("don't forget this...")


Or create SprintXxx functions to mix strings with other non-colorized strings:

    yellow := New(FgYellow).SprintFunc()
    red := New(FgRed).SprintFunc()

    fmt.Printf("this is a %s and this is %s.\n", yellow("warning"), red("error"))

    info := New(FgWhite, BgGreen).SprintFunc()
    fmt.Printf("this %s rocks!\n", info("package"))


Using with existing code is possible. Just use the Set() method to set the
standard output to the given parameters. That way a rewrite of an existing
code is not required.

    // Use handy standard colors.
    color.Set(color.FgYellow)

    fmt.Println("Existing text will be now in Yellow")
    fmt.Printf("This one %s\n", "too")

    color.Unset() // don't forget to unset

    // You can mix up parameters
    color.Set(color.FgMagenta, color.Bold)
    defer color.Unset() // use it in your function

    fmt.Println("All text will be now bold magenta.")
*/
package color
