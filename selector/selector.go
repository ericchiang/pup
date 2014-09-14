package selector

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"regexp"
	"strings"
)

// A CSS Selector
type Selector struct {
	Name  *regexp.Regexp
	Attrs map[string]*regexp.Regexp
}

type selectorField int

const (
	ClassField selectorField = iota
	IDField
	NameField
	AttrField
)

// Parse an attribute command to a key string and a regexp
func parseAttrField(command string) (attrKey string, matcher *regexp.Regexp,
	err error) {

	attrSplit := strings.Split(command, "=")
	matcherString := ""
	switch len(attrSplit) {
	case 1:
		attrKey = attrSplit[0]
		matcherString = ".*"
	case 2:
		attrKey = attrSplit[0]
		attrVal := attrSplit[1]
		if len(attrKey) == 0 {
			err = fmt.Errorf("No attribute key")
			return
		}
		attrKeyLen := len(attrKey)
		switch attrKey[attrKeyLen-1] {
		case '~':
			matcherString = fmt.Sprintf(`\b%s\b`, attrVal)
		case '$':
			matcherString = fmt.Sprintf("%s$", attrVal)
		case '^':
			matcherString = fmt.Sprintf("^%s", attrVal)
		case '*':
			matcherString = fmt.Sprintf("%s", attrVal)
		default:
			attrKeyLen++
			matcherString = fmt.Sprintf("^%s$", attrVal)
		}
		attrKey = attrKey[:attrKeyLen-1]
	default:
		err = fmt.Errorf("more than one '='")
		return
	}
	matcher, err = regexp.Compile(matcherString)
	return
}

// Set a field of this selector.
func (s *Selector) setFieldValue(f selectorField, v string) error {
	if v == "" {
		return nil
	}
	switch f {
	case ClassField:
		r, err := regexp.Compile(fmt.Sprintf(`\b%s\b`, v))
		if err != nil {
			return err
		}
		s.Attrs["class"] = r
	case IDField:
		r, err := regexp.Compile(fmt.Sprintf("^%s$", v))
		if err != nil {
			return err
		}
		s.Attrs["id"] = r
	case NameField:
		r, err := regexp.Compile(fmt.Sprintf("^%s$", v))
		if err != nil {
			return err
		}
		s.Name = r
	case AttrField:
		// Attribute fields are a little more complicated
		keystring, matcher, err := parseAttrField(v)
		if err != nil {
			return err
		}
		s.Attrs[keystring] = matcher
	}
	return nil
}

// Convert a string to a selector.
func NewSelector(s string) (*Selector, error) {
	attrs := map[string]*regexp.Regexp{}
	selector := &Selector{nil, attrs}
	nextField := NameField
	start := 0
	// Parse the selector character by character
	for i, c := range s {
		switch c {
		case '.':
			if nextField == AttrField {
				continue
			}
			err := selector.setFieldValue(nextField, s[start:i])
			if err != nil {
				return selector, err
			}
			nextField = ClassField
			start = i + 1
		case '#':
			if nextField == AttrField {
				continue
			}
			err := selector.setFieldValue(nextField, s[start:i])
			if err != nil {
				return selector, err
			}
			nextField = IDField
			start = i + 1
		case '[':
			err := selector.setFieldValue(nextField, s[start:i])
			if err != nil {
				return selector, err
			}
			nextField = AttrField
			start = i + 1
		case ']':
			if nextField != AttrField {
				return selector, fmt.Errorf(
					"']' must be preceeded by '['")
			}
			err := selector.setFieldValue(nextField, s[start:i])
			if err != nil {
				return selector, err
			}
			start = i + 1
		}
	}
	err := selector.setFieldValue(nextField, s[start:])
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
	matchedAttrs := []string{}
	for _, attr := range node.Attr {
		matcher, ok := sel.Attrs[attr.Key]
		if !ok {
			continue
		}
		if !matcher.MatchString(attr.Val) {
			return false
		}
		matchedAttrs = append(matchedAttrs, attr.Key)
	}
	for k := range sel.Attrs {
		attrMatched := false
		for _, attrKey := range matchedAttrs {
			if k == attrKey {
				attrMatched = true
			}
		}
		if !attrMatched {
			return false
		}
	}
	return true
}
