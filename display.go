package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

func init() {
	color.Output = colorable.NewColorableStdout()
}

type Displayer interface {
	Display([]*html.Node)
}

func ParseDisplayer(cmd string) error {
	attrRe := regexp.MustCompile(`attr\{([a-zA-Z\-]+)\}`)
	if cmd == "text{}" {
		pupDisplayer = TextDisplayer{}
	} else if cmd == "json{}" {
		pupDisplayer = JSONDisplayer{}
	} else if match := attrRe.FindAllStringSubmatch(cmd, -1); len(match) == 1 {
		pupDisplayer = AttrDisplayer{
			Attr: match[0][1],
		}
	} else {
		return fmt.Errorf("Unknown displayer")
	}
	return nil
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

var (
	// Colors
	tagColor     *color.Color = color.New(color.FgCyan)
	tokenColor                = color.New(color.FgCyan)
	attrKeyColor              = color.New(color.FgMagenta)
	quoteColor                = color.New(color.FgBlue)
	commentColor              = color.New(color.FgYellow)
)

type TreeDisplayer struct {
}

func (t TreeDisplayer) Display(nodes []*html.Node) {
	for _, node := range nodes {
		t.printNode(node, 0)
	}
}

// Print a node and all of it's children to `maxlevel`.
func (t TreeDisplayer) printNode(n *html.Node, level int) {
	switch n.Type {
	case html.TextNode:
		s := html.EscapeString(n.Data)
		s = strings.TrimSpace(s)
		if s != "" {
			t.printIndent(level)
			fmt.Println(s)
		}
	case html.ElementNode:
		t.printIndent(level)
		if pupPrintColor {
			tokenColor.Print("<")
			tagColor.Printf("%s", n.Data)
		} else {
			fmt.Printf("<%s", n.Data)
		}
		for _, a := range n.Attr {
			val := html.EscapeString(a.Val)
			if pupPrintColor {
				fmt.Print(" ")
				attrKeyColor.Printf("%s", a.Key)
				tokenColor.Print("=")
				quoteColor.Printf(`"%s"`, val)
			} else {
				fmt.Printf(` %s="%s"`, a.Key, val)
			}
		}
		if pupPrintColor {
			tokenColor.Println(">")
		} else {
			fmt.Println(">")
		}
		if !isVoidElement(n) {
			t.printChildren(n, level+1)
			t.printIndent(level)
			if pupPrintColor {
				tokenColor.Print("</")
				tagColor.Printf("%s", n.Data)
				tokenColor.Println(">")
			} else {
				fmt.Printf("</%s>\n", n.Data)
			}
		}
	case html.CommentNode:
		t.printIndent(level)
		if pupPrintColor {
			commentColor.Printf("<!--%s-->\n", n.Data)
		} else {
			fmt.Printf("<!--%s-->\n", n.Data)
		}
		t.printChildren(n, level)
	case html.DoctypeNode, html.DocumentNode:
		t.printChildren(n, level)
	}
}

func (t TreeDisplayer) printChildren(n *html.Node, level int) {
	if pupMaxPrintLevel > -1 {
		if level >= pupMaxPrintLevel {
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
		fmt.Print(pupIndentString)
	}
}

// Print the text of a node
type TextDisplayer struct{}

func (t TextDisplayer) Display(nodes []*html.Node) {
	for _, node := range nodes {
		if node.Type == html.TextNode {
			fmt.Println(node.Data)
		}
		children := []*html.Node{}
		child := node.FirstChild
		for child != nil {
			children = append(children, child)
			child = child.NextSibling
		}
		t.Display(children)
	}
}

// Print the attribute of a node
type AttrDisplayer struct {
	Attr string
}

func (a AttrDisplayer) Display(nodes []*html.Node) {
	for _, node := range nodes {
		attributes := node.Attr
		for _, attr := range attributes {
			if attr.Key == a.Attr {
				val := html.EscapeString(attr.Val)
				fmt.Printf("%s\n", val)
			}
		}
	}
}

// Print nodes as a JSON list
type JSONDisplayer struct{}

// returns a jsonifiable struct
func jsonify(node *html.Node) map[string]interface{} {
	vals := map[string]interface{}{}
	if len(node.Attr) > 0 {
		for _, attr := range node.Attr {
			vals[attr.Key] = html.EscapeString(attr.Val)
		}
	}
	vals["tag"] = node.DataAtom.String()
	children := []interface{}{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.ElementNode:
			children = append(children, jsonify(child))
		case html.TextNode:
			text := strings.TrimSpace(child.Data)
			if text != "" {
				// if there is already text we'll append it
				currText, ok := vals["text"]
				if ok {
					text = fmt.Sprintf("%s %s", currText, text)
				}
				vals["text"] = text
			}
		case html.CommentNode:
			comment := strings.TrimSpace(child.Data)
			currComment, ok := vals["comment"]
			if ok {
				comment = fmt.Sprintf("%s %s", currComment, comment)
			}
			vals["comment"] = comment
		}
	}
	if len(children) > 0 {
		vals["children"] = children
	}
	return vals
}

func (j JSONDisplayer) Display(nodes []*html.Node) {
	var data []byte
	var err error
	switch len(nodes) {
	case 1:
		node := nodes[0]
		if node.Type != html.DocumentNode {
			jsonNode := jsonify(nodes[0])
			data, err = json.MarshalIndent(&jsonNode, "", pupIndentString)
		} else {
			children := []*html.Node{}
			child := node.FirstChild
			for child != nil {
				children = append(children, child)
				child = child.NextSibling
			}
			j.Display(children)
		}
	default:
		jsonNodes := []map[string]interface{}{}
		for _, node := range nodes {
			jsonNodes = append(jsonNodes, jsonify(node))
		}
		data, err = json.MarshalIndent(&jsonNodes, "", pupIndentString)
	}
	if err != nil {
		panic("Could not jsonify nodes")
	}
	fmt.Printf("%s", data)
}
