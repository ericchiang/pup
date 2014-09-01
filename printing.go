package main

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"fmt"
	"github.com/fatih/color"
)

var (
	tagColor     *color.Color = color.New(color.FgYellow).Add(color.Bold)
	tokenColor                = color.New(color.FgCyan).Add(color.Bold)
	attrKeyColor              = color.New(color.FgRed)
	quoteColor                = color.New(color.FgBlue)
)

func printIndent(level int) {
	for ; level > 0; level-- {
		fmt.Print(indentString)
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
		if printColor {
			tokenColor.Print("<")
			tagColor.Printf("%s", n.Data)
		} else {
			fmt.Printf("<%s", n.Data)
		}
		for _, a := range n.Attr {
			if printColor {
				fmt.Print(" ")
				attrKeyColor.Printf("%s", a.Key)
				tokenColor.Print("=")
				quoteColor.Printf(`"%s"`, a.Val)
			} else {
				fmt.Printf(` %s="%s"`, a.Key, a.Val)
			}
		}
		if printColor {
			tokenColor.Println(">")
		} else {
			fmt.Print(">\n")
		}
		if !isVoidElement(n) {
			printChildren(n, level+1)
			printIndent(level)
			if printColor {
				tokenColor.Print("</")
				tagColor.Printf("%s", n.Data)
				tokenColor.Println(">")
			} else {
				fmt.Printf("</%s>\n", n.Data)
			}
		}
	case html.CommentNode, html.DoctypeNode, html.DocumentNode:
		printChildren(n, level)
	}
}
