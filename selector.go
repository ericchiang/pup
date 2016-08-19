package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"

	"golang.org/x/net/html"
)

type Selector interface {
	Match(node *html.Node) bool
}

type SelectorFunc func(nodes []*html.Node) []*html.Node

func Select(s Selector) SelectorFunc {
	// have to define first to be able to do recursion
	var selectChildren func(node *html.Node) []*html.Node
	selectChildren = func(node *html.Node) []*html.Node {
		selected := []*html.Node{}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if s.Match(child) {
				selected = append(selected, child)
			} else {
				selected = append(selected, selectChildren(child)...)
			}
		}
		return selected
	}
	return func(nodes []*html.Node) []*html.Node {
		selected := []*html.Node{}
		for _, node := range nodes {
			selected = append(selected, selectChildren(node)...)
		}
		return selected
	}
}

// Defined for the '>' selector
func SelectNextSibling(s Selector) SelectorFunc {
	return func(nodes []*html.Node) []*html.Node {
		selected := []*html.Node{}
		for _, node := range nodes {
			for ns := node.NextSibling; ns != nil; ns = ns.NextSibling {
				if ns.Type == html.ElementNode {
					if s.Match(ns) {
						selected = append(selected, ns)
					}
					break
				}
			}
		}
		return selected
	}
}

// Defined for the '+' selector
func SelectFromChildren(s Selector) SelectorFunc {
	return func(nodes []*html.Node) []*html.Node {
		selected := []*html.Node{}
		for _, node := range nodes {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				if s.Match(c) {
					selected = append(selected, c)
				}
			}
		}
		return selected
	}
}

type PseudoClass func(*html.Node) bool

type CSSSelector struct {
	Tag    string
	Attrs  map[string]*regexp.Regexp
	Pseudo PseudoClass
}

func (s CSSSelector) Match(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}
	if s.Tag != "" {
		if s.Tag != node.DataAtom.String() {
			return false
		}
	}
	for attrKey, matcher := range s.Attrs {
		matched := false
		for _, attr := range node.Attr {
			if attrKey == attr.Key {
				if !matcher.MatchString(attr.Val) {
					return false
				}
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	if s.Pseudo == nil {
		return true
	}
	return s.Pseudo(node)
}

// Parse a selector
// e.g. `div#my-button.btn[href^="http"]`
func ParseSelector(cmd string) (selector CSSSelector, err error) {
	selector = CSSSelector{
		Tag:    "",
		Attrs:  map[string]*regexp.Regexp{},
		Pseudo: nil,
	}
	var s scanner.Scanner
	s.Init(strings.NewReader(cmd))
	err = ParseTagMatcher(&selector, s)
	return
}

// Parse the initial tag
// e.g. `div`
func ParseTagMatcher(selector *CSSSelector, s scanner.Scanner) error {
	tag := bytes.NewBuffer([]byte{})
	defer func() {
		selector.Tag = tag.String()
	}()
	for {
		c := s.Next()
		switch c {
		case scanner.EOF:
			return nil
		case '.':
			return ParseClassMatcher(selector, s)
		case '#':
			return ParseIdMatcher(selector, s)
		case '[':
			return ParseAttrMatcher(selector, s)
		case ':':
			return ParsePseudo(selector, s)
		default:
			if _, err := tag.WriteRune(c); err != nil {
				return err
			}
		}
	}
}

// Parse a class matcher
// e.g. `.btn`
func ParseClassMatcher(selector *CSSSelector, s scanner.Scanner) error {
	var class bytes.Buffer
	defer func() {
		regexpStr := `(\A|\s)` + regexp.QuoteMeta(class.String()) + `(\s|\z)`
		selector.Attrs["class"] = regexp.MustCompile(regexpStr)
	}()
	for {
		c := s.Next()
		switch c {
		case scanner.EOF:
			return nil
		case '.':
			return ParseClassMatcher(selector, s)
		case '#':
			return ParseIdMatcher(selector, s)
		case '[':
			return ParseAttrMatcher(selector, s)
		case ':':
			return ParsePseudo(selector, s)
		default:
			if _, err := class.WriteRune(c); err != nil {
				return err
			}
		}
	}
}

// Parse an id matcher
// e.g. `#my-picture`
func ParseIdMatcher(selector *CSSSelector, s scanner.Scanner) error {
	var id bytes.Buffer
	defer func() {
		regexpStr := `^` + regexp.QuoteMeta(id.String()) + `$`
		selector.Attrs["id"] = regexp.MustCompile(regexpStr)
	}()
	for {
		c := s.Next()
		switch c {
		case scanner.EOF:
			return nil
		case '.':
			return ParseClassMatcher(selector, s)
		case '#':
			return ParseIdMatcher(selector, s)
		case '[':
			return ParseAttrMatcher(selector, s)
		case ':':
			return ParsePseudo(selector, s)
		default:
			if _, err := id.WriteRune(c); err != nil {
				return err
			}
		}
	}
}

// Parse an attribute matcher
// e.g. `[attr^="http"]`
func ParseAttrMatcher(selector *CSSSelector, s scanner.Scanner) error {
	var attrKey bytes.Buffer
	var attrVal bytes.Buffer
	hasMatchVal := false
	matchType := '='
	defer func() {
		if hasMatchVal {
			var regexpStr string
			switch matchType {
			case '=':
				regexpStr = `^` + regexp.QuoteMeta(attrVal.String()) + `$`
			case '*':
				regexpStr = regexp.QuoteMeta(attrVal.String())
			case '$':
				regexpStr = regexp.QuoteMeta(attrVal.String()) + `$`
			case '^':
				regexpStr = `^` + regexp.QuoteMeta(attrVal.String())
			case '~':
				regexpStr = `(\A|\s)` + regexp.QuoteMeta(attrVal.String()) + `(\s|\z)`
			}
			selector.Attrs[attrKey.String()] = regexp.MustCompile(regexpStr)
		} else {
			selector.Attrs[attrKey.String()] = regexp.MustCompile(`^.*$`)
		}
	}()
	// After reaching ']' proceed
	proceed := func() error {
		switch s.Next() {
		case scanner.EOF:
			return nil
		case '.':
			return ParseClassMatcher(selector, s)
		case '#':
			return ParseIdMatcher(selector, s)
		case '[':
			return ParseAttrMatcher(selector, s)
		case ':':
			return ParsePseudo(selector, s)
		default:
			return fmt.Errorf("Expected selector indicator after ']'")
		}
	}
	// Parse the attribute key matcher
	for !hasMatchVal {
		c := s.Next()
		switch c {
		case scanner.EOF:
			return fmt.Errorf("Unmatched open brace '['")
		case ']':
			// No attribute value matcher, proceed!
			return proceed()
		case '$', '^', '~', '*':
			matchType = c
			hasMatchVal = true
			if s.Next() != '=' {
				return fmt.Errorf("'%c' must be followed by a '='", matchType)
			}
		case '=':
			matchType = c
			hasMatchVal = true
		default:
			if _, err := attrKey.WriteRune(c); err != nil {
				return err
			}
		}
	}
	// figure out if the value is quoted
	c := s.Next()
	inQuote := false
	switch c {
	case scanner.EOF:
		return fmt.Errorf("Unmatched open brace '['")
	case ']':
		return proceed()
	case '"':
		inQuote = true
	default:
		if _, err := attrVal.WriteRune(c); err != nil {
			return err
		}
	}
	if inQuote {
		for {
			c := s.Next()
			switch c {
			case '\\':
				// consume another character
				if c = s.Next(); c == scanner.EOF {
					return fmt.Errorf("Unmatched open brace '['")
				}
			case '"':
				switch s.Next() {
				case ']':
					return proceed()
				default:
					return fmt.Errorf("Quote must end at ']'")
				}
			}
			if _, err := attrVal.WriteRune(c); err != nil {
				return err
			}
		}
	} else {
		for {
			c := s.Next()
			switch c {
			case scanner.EOF:
				return fmt.Errorf("Unmatched open brace '['")
			case ']':
				// No attribute value matcher, proceed!
				return proceed()
			}
			if _, err := attrVal.WriteRune(c); err != nil {
				return err
			}
		}
	}
}

// Parse the selector after ':'
func ParsePseudo(selector *CSSSelector, s scanner.Scanner) error {
	if selector.Pseudo != nil {
		return fmt.Errorf("Combined multiple pseudo classes")
	}
	var b bytes.Buffer
	for s.Peek() != scanner.EOF {
		if _, err := b.WriteRune(s.Next()); err != nil {
			return err
		}
	}
	cmd := b.String()
	var err error
	switch {
	case cmd == "empty":
		selector.Pseudo = func(n *html.Node) bool {
			return n.FirstChild == nil
		}
	case cmd == "first-child":
		selector.Pseudo = firstChildPseudo
	case cmd == "last-child":
		selector.Pseudo = lastChildPseudo
	case cmd == "only-child":
		selector.Pseudo = func(n *html.Node) bool {
			return firstChildPseudo(n) && lastChildPseudo(n)
		}
	case cmd == "first-of-type":
		selector.Pseudo = firstOfTypePseudo
	case cmd == "last-of-type":
		selector.Pseudo = lastOfTypePseudo
	case cmd == "only-of-type":
		selector.Pseudo = func(n *html.Node) bool {
			return firstOfTypePseudo(n) && lastOfTypePseudo(n)
		}
	case strings.HasPrefix(cmd, "contains("):
		selector.Pseudo, err = parseContainsPseudo(cmd[len("contains("):])
		if err != nil {
			return err
		}
	case strings.HasPrefix(cmd, "nth-child("),
		strings.HasPrefix(cmd, "nth-last-child("),
		strings.HasPrefix(cmd, "nth-last-of-type("),
		strings.HasPrefix(cmd, "nth-of-type("):
		if selector.Pseudo, err = parseNthPseudo(cmd); err != nil {
			return err
		}
	case strings.HasPrefix(cmd, "not("):
		if selector.Pseudo, err = parseNotPseudo(cmd[len("not("):]); err != nil {
			return err
		}
	case strings.HasPrefix(cmd, "parent-of("):
		if selector.Pseudo, err = parseParentOfPseudo(cmd[len("parent-of("):]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s not a valid pseudo class", cmd)
	}
	return nil
}

// :first-of-child
func firstChildPseudo(n *html.Node) bool {
	for c := n.PrevSibling; c != nil; c = c.PrevSibling {
		if c.Type == html.ElementNode {
			return false
		}
	}
	return true
}

// :last-of-child
func lastChildPseudo(n *html.Node) bool {
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			return false
		}
	}
	return true
}

// :first-of-type
func firstOfTypePseudo(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}
	for n := node.PrevSibling; n != nil; n = n.PrevSibling {
		if n.DataAtom == node.DataAtom {
			return false
		}
	}
	return true
}

// :last-of-type
func lastOfTypePseudo(node *html.Node) bool {
	if node.Type != html.ElementNode {
		return false
	}
	for n := node.NextSibling; n != nil; n = n.NextSibling {
		if n.DataAtom == node.DataAtom {
			return false
		}
	}
	return true
}

func parseNthPseudo(cmd string) (PseudoClass, error) {
	i := strings.IndexRune(cmd, '(')
	if i < 0 {
		// really, we should never get here
		return nil, fmt.Errorf("Fatal error, '%s' does not contain a '('", cmd)
	}
	pseudoName := cmd[:i]
	// Figure out how the counting function works
	var countNth func(*html.Node) int
	switch pseudoName {
	case "nth-child":
		countNth = func(n *html.Node) int {
			nth := 1
			for sib := n.PrevSibling; sib != nil; sib = sib.PrevSibling {
				if sib.Type == html.ElementNode {
					nth++
				}
			}
			return nth
		}
	case "nth-of-type":
		countNth = func(n *html.Node) int {
			nth := 1
			for sib := n.PrevSibling; sib != nil; sib = sib.PrevSibling {
				if sib.Type == html.ElementNode && sib.DataAtom == n.DataAtom {
					nth++
				}
			}
			return nth
		}
	case "nth-last-child":
		countNth = func(n *html.Node) int {
			nth := 1
			for sib := n.NextSibling; sib != nil; sib = sib.NextSibling {
				if sib.Type == html.ElementNode {
					nth++
				}
			}
			return nth
		}
	case "nth-last-of-type":
		countNth = func(n *html.Node) int {
			nth := 1
			for sib := n.NextSibling; sib != nil; sib = sib.NextSibling {
				if sib.Type == html.ElementNode && sib.DataAtom == n.DataAtom {
					nth++
				}
			}
			return nth
		}
	default:
		return nil, fmt.Errorf("Unrecognized pseudo '%s'", pseudoName)
	}

	nthString := cmd[i+1:]
	i = strings.IndexRune(nthString, ')')
	if i < 0 {
		return nil, fmt.Errorf("Unmatched '(' for pseudo class %s", pseudoName)
	} else if i != len(nthString)-1 {
		return nil, fmt.Errorf("%s(n) must end selector", pseudoName)
	}
	number := nthString[:i]

	// Check if the number is 'odd' or 'even'
	oddOrEven := -1
	switch number {
	case "odd":
		oddOrEven = 1
	case "even":
		oddOrEven = 0
	}
	if oddOrEven > -1 {
		return func(n *html.Node) bool {
			return n.Type == html.ElementNode && countNth(n)%2 == oddOrEven
		}, nil
	}
	// Check against '3n+4' pattern
	r := regexp.MustCompile(`([0-9]+)n[ ]?\+[ ]?([0-9])`)
	subMatch := r.FindAllStringSubmatch(number, -1)
	if len(subMatch) == 1 && len(subMatch[0]) == 3 {
		cycle, _ := strconv.Atoi(subMatch[0][1])
		offset, _ := strconv.Atoi(subMatch[0][2])
		return func(n *html.Node) bool {
			return n.Type == html.ElementNode && countNth(n)%cycle == offset
		}, nil
	}
	// check against 'n+2' pattern
	r = regexp.MustCompile(`n[ ]?\+[ ]?([0-9])`)
	subMatch = r.FindAllStringSubmatch(number, -1)
	if len(subMatch) == 1 && len(subMatch[0]) == 2 {
		offset, _ := strconv.Atoi(subMatch[0][1])
		return func(n *html.Node) bool {
			return n.Type == html.ElementNode && countNth(n) >= offset
		}, nil
	}
	// the only other option is a numeric value
	nth, err := strconv.Atoi(nthString[:i])
	if err != nil {
		return nil, err
	} else if nth <= 0 {
		return nil, fmt.Errorf("Argument to '%s' must be greater than 0", pseudoName)
	}
	return func(n *html.Node) bool {
		return n.Type == html.ElementNode && countNth(n) == nth
	}, nil
}

// Parse a :contains("") selector
// expects the input to be everything after the open parenthesis
// e.g. for `contains("Help")` the argument would be `"Help")`
func parseContainsPseudo(cmd string) (PseudoClass, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(cmd))
	switch s.Next() {
	case '"':
	default:
		return nil, fmt.Errorf("Malformed 'contains(\"\")' selector")
	}
	textToContain := bytes.NewBuffer([]byte{})
	for {
		r := s.Next()
		switch r {
		case '"':
			// ')' then EOF must follow '"'
			if s.Next() != ')' {
				return nil, fmt.Errorf("Malformed 'contains(\"\")' selector")
			}
			if s.Next() != scanner.EOF {
				return nil, fmt.Errorf("'contains(\"\")' must end selector")
			}
			text := textToContain.String()
			contains := func(node *html.Node) bool {
				for c := node.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						if strings.Contains(c.Data, text) {
							return true
						}
					}
				}
				return false
			}
			return contains, nil
		case '\\':
			s.Next()
		case scanner.EOF:
			return nil, fmt.Errorf("Malformed 'contains(\"\")' selector")
		default:
			if _, err := textToContain.WriteRune(r); err != nil {
				return nil, err
			}
		}
	}
}

// Parse a :not(selector) selector
// expects the input to be everything after the open parenthesis
// e.g. for `not(div#id)` the argument would be `div#id)`
func parseNotPseudo(cmd string) (PseudoClass, error) {
	if len(cmd) < 2 {
		return nil, fmt.Errorf("malformed ':not' selector")
	}
	endQuote, cmd := cmd[len(cmd)-1], cmd[:len(cmd)-1]
	selector, err := ParseSelector(cmd)
	if err != nil {
		return nil, err
	}
	if endQuote != ')' {
		return nil, fmt.Errorf("unmatched '('")
	}
	return func(n *html.Node) bool {
		return !selector.Match(n)
	}, nil
}

// Parse a :parent-of(selector) selector
// expects the input to be everything after the open parenthesis
// e.g. for `parent-of(div#id)` the argument would be `div#id)`
func parseParentOfPseudo(cmd string) (PseudoClass, error) {
	if len(cmd) < 2 {
		return nil, fmt.Errorf("malformed ':parent-of' selector")
	}
	endQuote, cmd := cmd[len(cmd)-1], cmd[:len(cmd)-1]
	selector, err := ParseSelector(cmd)
	if err != nil {
		return nil, err
	}
	if endQuote != ')' {
		return nil, fmt.Errorf("unmatched '('")
	}
	return func(n *html.Node) bool {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && selector.Match(c) {
				return true
			}
		}
		return false
	}, nil
}
