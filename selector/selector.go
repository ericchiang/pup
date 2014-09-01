package selector

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"regexp"
	"strings"
)

// A CSS Selector
type Selector struct {
	Class, ID, Name *regexp.Regexp
	Attrs           map[string]*regexp.Regexp
}

type selectorField string

const (
	Class selectorField = "class"
	ID    selectorField = "id"
	Name  selectorField = "name"
)

// Set a field of this selector.
func (s *Selector) setFieldValue(a selectorField, v string) error {
	if v == "" {
		return nil
	}
	// wildcards become '.*'
	v = strings.Replace(v, "*", ".*", -1)
	r, err := regexp.Compile(fmt.Sprintf("^%s$", v))
	if err != nil {
		return err
	}
	switch a {
	case Class:
		s.Class = r
	case ID:
		s.ID = r
	case Name:
		s.Name = r
	}
	return nil
}

// Convert a string to a selector.
func NewSelector(s string) (*Selector, error) {
	attrs := map[string]*regexp.Regexp{}
	selector := &Selector{nil, nil, nil, attrs}
	nextAttr := Name
	start := 0
	for i, c := range s {
		switch c {
		case '.':
			err := selector.setFieldValue(nextAttr, s[start:i])
			if err != nil {
				return selector, err
			}
			nextAttr = Class
			start = i + 1
		case '#':
			err := selector.setFieldValue(nextAttr, s[start:i])
			if err != nil {
				return selector, err
			}
			nextAttr = ID
			start = i + 1
		}
	}
	err := selector.setFieldValue(nextAttr, s[start:])
	if err != nil {
		return selector, err
	}
	return selector, nil
}

// Find all nodes which match a selector.
func (sel *Selector) FindAllChildren(node *html.Node) []*html.Node {
	selected := []*html.Node{}
	child := node.FirstChild
	for child != nil {
		childSelected := sel.FindAll(child)
		selected = append(selected, childSelected...)
		child = child.NextSibling
	}
	return selected
}

// Find all nodes which match a selector. May return itself.
func (sel *Selector) FindAll(node *html.Node) []*html.Node {
	selected := []*html.Node{}
	if sel.Match(node) {
		return []*html.Node{node}
	}
	child := node.FirstChild
	for child != nil {
		childSelected := sel.FindAll(child)
		selected = append(selected, childSelected...)
		child = child.NextSibling
	}
	return selected
}

// Does this selector match a given node?
func (sel *Selector) Match(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}
	if sel.Name != nil {
		if !sel.Name.MatchString(strings.ToLower(node.Data)) {
			return false
		}
	}
	classMatched := sel.Class == nil
	idMatched := sel.ID == nil
	for _, attr := range node.Attr {
		switch attr.Key {
		case "class":
			if !classMatched {
				if !sel.Class.MatchString(attr.Val) {
					return false
				} else {
					classMatched = true
				}
			}
		case "id":
			if !idMatched {
				if !sel.ID.MatchString(attr.Val) {
					return false
				} else {
					idMatched = true
				}
			}
		}
	}
	return classMatched && idMatched
}
