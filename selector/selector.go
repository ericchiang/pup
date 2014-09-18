package selector

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// A CSS Selector
type BasicSelector struct {
	Name  *regexp.Regexp
	Attrs map[string]*regexp.Regexp
}

type Selector interface {
	Select(nodes []*html.Node) []*html.Node
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
func (s *BasicSelector) setFieldValue(f selectorField, v string) error {
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
func NewSelector(s string) (Selector, error) {
	// A very simple test for a selector function
	if strings.Contains(s, "{") {
		return parseSelectorFunc(s)
	}

	// Otherwise let's evaluate a basic selector
	attrs := map[string]*regexp.Regexp{}
	selector := BasicSelector{nil, attrs}
	nextField := NameField
	start := 0
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

func (sel BasicSelector) Select(nodes []*html.Node) []*html.Node {
	selected := []*html.Node{}
	for _, node := range nodes {
		selected = append(selected, sel.FindAllChildren(node)...)
	}
	return selected
}

// Find all nodes which match a selector.
func (sel BasicSelector) FindAllChildren(node *html.Node) []*html.Node {
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
func (sel BasicSelector) FindAll(node *html.Node) []*html.Node {
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
func (sel BasicSelector) Match(node *html.Node) bool {
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

type SliceSelector struct {
	Start      int
	LimitStart bool
	End        int
	LimitEnd   bool
	By         int
}

func (sel SliceSelector) Select(nodes []*html.Node) []*html.Node {
	var start, end, by int
	selected := []*html.Node{}
	nNodes := len(nodes)
	switch {
	case !sel.LimitStart:
		start = 0
	case sel.Start < 0:
		start = (nNodes + 1) + sel.Start
	default:
		start = sel.Start
	}
	switch {
	case !sel.LimitEnd:
		end = nNodes
	case sel.End < 0:
		end = (nNodes + 1) + sel.End
	default:
		end = sel.End
	}
	by = sel.By
	if by == 0 {
		return selected
	}
	if by > 0 {
		for i := start; i < nNodes && i < end; i = i + by {
			selected = append(selected, nodes[i])
		}
	} else {
		for i := end - 1; i > 0 && i >= start; i = i + by {
			selected = append(selected, nodes[i])
		}
	}
	return selected
}

// expects input to be the slice only, e.g. "9:4:-1"
func parseSliceSelector(s string) (sel SliceSelector, err error) {
	sel = SliceSelector{
		Start:      0,
		End:        0,
		By:         1,
		LimitStart: false,
		LimitEnd:   false,
	}
	split := strings.Split(s, ":")
	n := len(split)
	if n > 3 {
		err = fmt.Errorf("too many slices")
		return
	}
	var value int
	if split[0] != "" {
		value, err = strconv.Atoi(split[0])
		if err != nil {
			return
		}
		sel.Start = value
		sel.LimitStart = true
	}
	if n == 1 {
		sel.End = sel.Start + 1
		sel.LimitEnd = true
		return
	}
	if split[1] != "" {
		value, err = strconv.Atoi(split[1])
		if err != nil {
			return
		}
		sel.End = value
		sel.LimitEnd = true
	}
	if n == 2 {
		return
	}
	if split[2] != "" {
		value, err = strconv.Atoi(split[2])
		if err != nil {
			return
		}
		sel.By = value
	}
	return
}

func parseSelectorFunc(s string) (Selector, error) {
	switch {
	case strings.HasPrefix(s, "{"):
		if !strings.HasSuffix(s, "}") {
			return nil, fmt.Errorf(
				"slice func must end with a '}'")
		}
		s = strings.TrimPrefix(s, "{")
		s = strings.TrimSuffix(s, "}")
		return parseSliceSelector(s)
	}
	return nil, fmt.Errorf("%s is an invalid function", s)
}
