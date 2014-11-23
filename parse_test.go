package main

import (
	"testing"
)

type parseCmdTest struct {
	input string
	split []string
	ok    bool
}

var parseCmdTests = []parseCmdTest{
	parseCmdTest{`w1 w2`, []string{`w1`, `w2`}, true},
	parseCmdTest{`w1 w2 w3`, []string{`w1`, `w2`, `w3`}, true},
	parseCmdTest{`w1 'w2 w3'`, []string{`w1`, `'w2 w3'`}, true},
	parseCmdTest{`w1 "w2 w3"`, []string{`w1`, `"w2 w3"`}, true},
	parseCmdTest{`w1   "w2 w3"`, []string{`w1`, `"w2 w3"`}, true},
	parseCmdTest{`w1   'w2 w3'`, []string{`w1`, `'w2 w3'`}, true},
	parseCmdTest{`w1"w2 w3"`, []string{`w1"w2 w3"`}, true},
	parseCmdTest{`w1'w2 w3'`, []string{`w1'w2 w3'`}, true},
	parseCmdTest{`w1"w2 'w3"`, []string{`w1"w2 'w3"`}, true},
	parseCmdTest{`w1'w2 "w3'`, []string{`w1'w2 "w3'`}, true},
	parseCmdTest{`"w1 w2" "w3"`, []string{`"w1 w2"`, `"w3"`}, true},
	parseCmdTest{`'w1 w2' "w3"`, []string{`'w1 w2'`, `"w3"`}, true},
	parseCmdTest{`'w1 \'w2' "w3"`, []string{`'w1 \'w2'`, `"w3"`}, true},
	parseCmdTest{`'w1 \'w2 "w3"`, []string{}, false},
	parseCmdTest{`w1 'w2 w3'"`, []string{}, false},
	parseCmdTest{`w1 "w2 w3"'`, []string{}, false},
	parseCmdTest{`w1 '  "w2 w3"`, []string{}, false},
	parseCmdTest{`w1 "  'w2 w3'`, []string{}, false},
	parseCmdTest{`w1"w2 w3""`, []string{}, false},
	parseCmdTest{`w1'w2 w3''`, []string{}, false},
	parseCmdTest{`w1"w2 'w3""`, []string{}, false},
	parseCmdTest{`w1'w2 "w3''`, []string{}, false},
	parseCmdTest{`"w1 w2" "w3"'`, []string{}, false},
	parseCmdTest{`'w1 w2' "w3"'`, []string{}, false},
	parseCmdTest{`w1,"w2 w3"`, []string{`w1`, `,`, `"w2 w3"`}, true},
	parseCmdTest{`w1,'w2 w3'`, []string{`w1`, `,`, `'w2 w3'`}, true},
	parseCmdTest{`w1  ,  "w2 w3"`, []string{`w1`, `,`, `"w2 w3"`}, true},
	parseCmdTest{`w1  ,  'w2 w3'`, []string{`w1`, `,`, `'w2 w3'`}, true},
	parseCmdTest{`w1,  "w2 w3"`, []string{`w1`, `,`, `"w2 w3"`}, true},
	parseCmdTest{`w1,  'w2 w3'`, []string{`w1`, `,`, `'w2 w3'`}, true},
	parseCmdTest{`w1  ,"w2 w3"`, []string{`w1`, `,`, `"w2 w3"`}, true},
	parseCmdTest{`w1  ,'w2 w3'`, []string{`w1`, `,`, `'w2 w3'`}, true},
	parseCmdTest{`w1"w2, w3"`, []string{`w1"w2, w3"`}, true},
	parseCmdTest{`w1'w2, w3'`, []string{`w1'w2, w3'`}, true},
	parseCmdTest{`w1"w2, 'w3"`, []string{`w1"w2, 'w3"`}, true},
	parseCmdTest{`w1'w2, "w3'`, []string{`w1'w2, "w3'`}, true},
	parseCmdTest{`"w1, w2" "w3"`, []string{`"w1, w2"`, `"w3"`}, true},
	parseCmdTest{`'w1, w2' "w3"`, []string{`'w1, w2'`, `"w3"`}, true},
	parseCmdTest{`'w1, \'w2' "w3"`, []string{`'w1, \'w2'`, `"w3"`}, true},
	parseCmdTest{`h1, .article-teaser, .article-content`, []string{
		`h1`, `,`, `.article-teaser`, `,`, `.article-content`,
	}, true},
	parseCmdTest{`h1 ,.article-teaser ,.article-content`, []string{
		`h1`, `,`, `.article-teaser`, `,`, `.article-content`,
	}, true},
	parseCmdTest{`h1 , .article-teaser , .article-content`, []string{
		`h1`, `,`, `.article-teaser`, `,`, `.article-content`,
	}, true},
}

func sliceEq(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func TestParseCommands(t *testing.T) {
	for _, test := range parseCmdTests {
		parsed, err := ParseCommands(test.input)
		if test.ok != (err == nil) {
			t.Errorf("`%s`: should have cause error? %v", test.input, !test.ok)
		} else if !sliceEq(test.split, parsed) {
			t.Errorf("`%s`: `%s`: `%s`", test.input, test.split, parsed)
		}
	}
}
