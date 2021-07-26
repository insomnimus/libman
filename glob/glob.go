package glob

import (
	"fmt"
	"regexp"
	"strings"
)

type Regexp = regexp.Regexp

func Compile(pattern string) (*Regexp, error) {
	c := compiler{
		expr: []rune(pattern),
	}
	c.read()
	return regexp.Compile(c.compile())
}

type compiler struct {
	ch           rune
	expr         []rune
	pos, readpos int
}

func (c *compiler) read() {
	if c.readpos >= len(c.expr) {
		c.ch = 0
	} else {
		c.ch = c.expr[c.readpos]
	}
	c.pos = c.readpos
	c.readpos++
}

func (c *compiler) peek() rune {
	if c.readpos >= len(c.expr) {
		return 0
	}
	return c.expr[c.readpos]
}

func (c *compiler) compile() string {
	var buf strings.Builder
	for c.ch != 0 {
		switch c.ch {
		case ' ':
			for c.peek() == ' ' {
				c.read()
			}
			buf.WriteString("\\s+?")
		case '*':
			for c.peek() == '*' {
				c.read()
			}
			buf.WriteString(".*")
		case '\\':
			if isRegSpecial(c.peek()) {
				buf.WriteRune('\\')
				c.read()
				buf.WriteRune(c.ch)
			} else {
				buf.WriteString("\\\\")
			}
		default:
			if isRegSpecial(c.ch) {
				buf.WriteRune('\\')
				buf.WriteRune(c.ch)
			} else {
				buf.WriteRune(c.ch)
			}
		}
		c.read()
	}
	return fmt.Sprintf("(?i)^%s$", &buf)
}

func isRegSpecial(c rune) bool {
	for _, ch := range `^$\.,*+-[]{}():?` {
		if c == ch {
			return true
		}
	}
	return false
}
