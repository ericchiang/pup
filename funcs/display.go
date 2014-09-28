package funcs

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"regexp"
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

var (
	// Display function helpers
	displayMatcher  *regexp.Regexp = regexp.MustCompile(`\{[^\}]*\}$`)
	textFuncMatcher                = regexp.MustCompile(`^text\{\}$`)
	attrFuncMatcher                = regexp.MustCompile(`^attr\{([^\}]*)\}$`)
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
	}
	return nil, fmt.Errorf("Not a display function")
}
