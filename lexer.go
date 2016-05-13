package BetterCobol

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type itemType int

type item struct {
	typ itemType
	val string
}

type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan item
}

type stateFn func(*lexer) stateFn

type searchField struct {
	search string
	result stateFn
}

const (
	itemEOF itemType = iota
	itemCopyBook
	itemDataDivision
	itemProcedureDivision
)

const eof = -1

func (l *lexer) run() {
	for state := lexSearch; state != nil; {
		state = state(l)
	}

	close(l.items)
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) searchFieldList(fields []searchField) stateFn {
	for _, f := range fields {
		if strings.HasPrefix(l.input[l.pos:], f.search) {
			return f.result
		}
	}
	return nil
}

func (l *lexer) acceptWhiteSpace() {
	for unicode.IsSpace(l.next()) {
	}
	l.backup()
}

func lex(name, input string, bufferSize int) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item, bufferSize),
	}

	go l.run()
	return l, l.items
}
