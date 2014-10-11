package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"
)

type Displayer interface {
	Display(nodes []*html.Node)
}

type TextDisplayer struct {
}

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

type JSONDisplayer struct {
}

// returns a jsonifiable struct
func jsonify(node *html.Node) map[string]interface{} {
	vals := map[string]interface{}{}
	if len(node.Attr) > 0 {
		attrs := map[string]string{}
		for _, attr := range node.Attr {
			attrs[attr.Key] = html.EscapeString(attr.Val)
		}
		vals["attrs"] = attrs
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
		}
	}
	return vals
}

func (j JSONDisplayer) Display(nodes []*html.Node) {
	var data []byte
	var err error
	switch len(nodes) {
	case 1:
		jsonNode := jsonify(nodes[0])
		data, err = json.MarshalIndent(&jsonNode, "", indentString)
	default:
		jsonNodes := []map[string]interface{}{}
		for _, node := range nodes {
			jsonNodes = append(jsonNodes, jsonify(node))
		}
		data, err = json.MarshalIndent(&jsonNodes, "", indentString)
	}
	if err != nil {
		panic("Could not jsonify nodes")
	}
	fmt.Printf("%s\n", data)
}

var (
	// Display function helpers
	displayMatcher  *regexp.Regexp = regexp.MustCompile(`\{[^\}]*\}$`)
	textFuncMatcher                = regexp.MustCompile(`^text\{\}$`)
	attrFuncMatcher                = regexp.MustCompile(`^attr\{([^\}]*)\}$`)
	jsonFuncMatcher                = regexp.MustCompile(`^json\{([^\}]*)\}$`)
)

func NewDisplayFunc(text string) (Displayer, error) {
	if !displayMatcher.MatchString(text) {
		return nil, fmt.Errorf("Not a display function")
	}
	switch {
	case textFuncMatcher.MatchString(text):
		return TextDisplayer{}, nil
	case attrFuncMatcher.MatchString(text):
		matches := attrFuncMatcher.FindStringSubmatch(text)
		if len(matches) != 2 {
			return nil, fmt.Errorf("")
		} else {
			return AttrDisplayer{matches[1]}, nil
		}
	case jsonFuncMatcher.MatchString(text):
		return JSONDisplayer{}, nil
	}
	return nil, fmt.Errorf("Not a display function")
}
