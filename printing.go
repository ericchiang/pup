package main

import (
	"fmt"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Colors
	tagColor     *color.Color = color.New(color.FgCyan)
	tokenColor                = color.New(color.FgCyan)
	attrKeyColor              = color.New(color.FgMagenta)
	quoteColor                = color.New(color.FgBlue)
	commentColor              = color.New(color.FgYellow)
)

func init() {
	color.Output = colorable.NewColorableStdout()
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

type TreeDisplayer struct {
	IndentString string
}

func (t TreeDisplayer) Display(nodes []*html.Node) {
	for _, node := range nodes {
		t.printNode(node, 0)
	}
}

func (t TreeDisplayer) printChildren(n *html.Node, level int) {
	if maxPrintLevel > -1 {
		if level >= maxPrintLevel {
			t.printIndent(level)
			fmt.Println("...")
			return
		}
	}
	child := n.FirstChild
	for child != nil {
		t.printNode(child, level)
		child = child.NextSibling
	}
}

func (t TreeDisplayer) printIndent(level int) {
	for ; level > 0; level-- {
		fmt.Print(indentString)
	}
}

// Print a node and all of it's children to `maxlevel`.
func (t TreeDisplayer) printNode(n *html.Node, level int) {
	switch n.Type {
	case html.TextNode:
		s := html.EscapeString(n.Data)
		s = strings.TrimSpace(s)
		if s != "" {
			t.printIndent(level + 1)
			fmt.Println(s)
		}
	case html.ElementNode:
		t.printIndent(level)
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
				val := html.EscapeString(a.Val)
				quoteColor.Printf(`"%s"`, val)
			} else {
				val := html.EscapeString(a.Val)
				fmt.Printf(` %s="%s"`, a.Key, val)
			}
		}
		if printColor {
			tokenColor.Println(">")
		} else {
			fmt.Print(">\n")
		}
		if !isVoidElement(n) {
			t.printChildren(n, level+1)
			t.printIndent(level)
			if printColor {
				tokenColor.Print("</")
				tagColor.Printf("%s", n.Data)
				tokenColor.Println(">")
			} else {
				fmt.Printf("</%s>\n", n.Data)
			}
		}
	case html.CommentNode:
		t.printIndent(level)
		if printColor {
			commentColor.Printf("<!--%s-->\n", n.Data)
		} else {
			fmt.Printf("<!--%s-->\n", n.Data)
		}
		t.printChildren(n, level)
	case html.DoctypeNode, html.DocumentNode:
		t.printChildren(n, level)
	}
}
