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

var (
	// Display function helpers
	displayMatcher  *regexp.Regexp = regexp.MustCompile(`\{[^\}]*\}$`)
	textFuncMatcher                = regexp.MustCompile(`^text\{\}$`)
	attrFuncMatcher                = regexp.MustCompile(`^attr\{[^\}]*\}$`)
)

func NewDisplayFunc(text string) (Displayer, error) {
	if !displayMatcher.MatchString(text) {
		return nil, fmt.Errorf("Not a display function")
	}
	switch {
	case textFuncMatcher.MatchString(text):
		return TextDisplayer{}, nil
	case attrFuncMatcher.MatchString(text):
		return nil, fmt.Errorf("attr")
	}
	return nil, fmt.Errorf("Not a display function")
}
