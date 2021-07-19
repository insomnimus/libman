package util

import (
	"testing"
)

func TestHasPrefixFold(t *testing.T) {
	tests := []struct {
		x, y string
	}{
		{"asdf", ""},
		{"asdf", "AsDf"},
		{"LMAO 123", "lmA"},
		{"ĞĞĞAfs", "ğğ"},
	}

	for _, test := range tests {
		if !HasPrefixFold(test.x, test.y) {
			t.Errorf("hasPrefixFold returned false\nleft: %s\nright: %s\n", test.x, test.y)
		}
	}
}

func TestSplitCmd(t *testing.T) {
	tests := []struct {
		s           string
		left, right string
	}{
		{"echo hi", "echo", "hi"},
		{"echo    hi  lmao", "echo", "hi  lmao"},
		{"echo  hi    lmao  ", "echo", "hi    lmao"},
	}

	for _, test := range tests {
		left, right := SplitCmd(test.s)
		if test.left != left {
			t.Errorf("splitCmd failed:\ntext: %q\nleft: %q\nright: %q\n\ngot:\nleft: %q\nright: %q\n", test.s, test.left, test.right, left, right)
		}
	}
}
