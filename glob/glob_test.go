package glob

import (
	"fmt"
	"testing"
)

func TestCompile(t *testing.T) {
	cases := []struct {
		x, y string
	}{
		{"abcd*", "abcd.*"},
		{"haha??.really", "haha\\?\\?\\.really"},
		{"**bang**", ".*bang.*"},
		{"nice, dude", "nice\\,\\s+?dude"},
	}

	for i, test := range cases {
		c := compiler{expr: []rune(test.x)}
		c.read()
		got := c.compile()
		expected := fmt.Sprintf("(?i)^%s$", test.y)
		if got != expected {
			t.Errorf("failed test #%d:\nexpected %s\ngot %s\n", i, expected, got)
		}
	}
}
