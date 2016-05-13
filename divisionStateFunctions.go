package BetterCobol

import (
	"strings"
	"unicode"
)

func lexDataStart(l *lexer) stateFn {
	l.pos += len(dataStart)
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case strings.HasPrefix(l.input[l.pos:], division):
			return lexDataDivision
		}
	}
}

func lexDataDivision(l *lexer) stateFn {
	l.pos += len(division)
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case r == '.':
			return lexDataDivisionContent
		case unicode.IsSpace(r):
			l.ignore()
		default:
			return cobolError
		}
	}
}

func lexDataDivisionContent(l *lexer) stateFn {
	l.ignore()
	for {
		r := l.next()
		switch {
		case r == eof:
			l.backup()
			l.emit(itemDataDivision)
			return lexSearch
		case strings.HasPrefix(l.input[l.pos:], procedureStart):
			l.emit(itemDataDivision)
			return lexSearch
		}
	}
}

func lexProcedureStart(l *lexer) stateFn {
	l.pos += len(procedureStart)
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case strings.HasPrefix(l.input[l.pos:], division):
			return lexProcedureDivision
		}
	}
}

func lexProcedureDivision(l *lexer) stateFn {
	l.pos += len(division)
	l.ignore()
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case r == '.':
			return lexProcedureDivisionContent
		case strings.HasPrefix(l.input[l.pos:], using):
			return lexProcedureDivisionUsing
		case unicode.IsSpace(r):
			l.ignore()
		default:
			return cobolError
		}
	}
}

func lexProcedureDivisionUsing(l *lexer) stateFn {
	l.pos += len(division)
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case r == '.':
			return lexProcedureDivisionContent
		}
	}
}

func lexProcedureDivisionContent(l *lexer) stateFn {
	l.ignore()
	for {
		r := l.next()
		switch {
		case r == eof:
			l.backup()
			l.emit(itemProcedureDivision)
			return lexSearch
		}
	}
}
