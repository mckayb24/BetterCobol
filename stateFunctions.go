package BetterCobol

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	execStart      = "EXEC"
	execSQL        = "SQL"
	execCopyBook   = "INCLUDE"
	execEnd        = "END-EXEC."
	dataStart      = "\n       DATA"
	division       = "DIVISION"
	procedureStart = "\n       PROCEDURE"
	using          = "USING"
)

func lexSearch(l *lexer) stateFn {
	items := []searchField{
		{execStart, lexExec},
		{dataStart, lexDataStart},
		{procedureStart, lexProcedureStart},
	}
	for {
		if nextState := l.searchFieldList(items); nextState != nil {
			if l.pos > l.start {
				l.ignore()
			}
			return nextState // Next state.
		}
		r := l.next()
		if r == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.ignore()
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

func lexExec(l *lexer) stateFn {
	l.pos += len(execStart)
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case strings.HasPrefix(l.input[l.pos:], execSQL):
			return lexExecSQL
		case strings.HasPrefix(l.input[l.pos:], execEnd):
			return lexSearch
		default:
			l.ignore()
		}

	}
	return nil
}

func lexExecSQL(l *lexer) stateFn {
	l.pos += len(execSQL)
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case strings.HasPrefix(l.input[l.pos:], execCopyBook):
			return lexCopyBook
		case strings.HasPrefix(l.input[l.pos:], execEnd):
			return lexSearch
		default:
			l.ignore()
		}
	}
	return nil
}

func lexCopyBook(l *lexer) stateFn {
	l.pos += len(execCopyBook)
	l.acceptWhiteSpace()
	l.ignore()
	for {
		r := l.next()
		switch {
		case r == eof:
			return cobolError
		case unicode.IsSpace(r):
			l.backup()
			l.emit(itemCopyBook)
			return lexExecEnd
		}
	}
}

func lexExecEnd(l *lexer) stateFn {
	for {
		l.next()
		if strings.HasPrefix(l.input[l.pos:], execEnd) {
			l.pos += len(execEnd)
			l.ignore()
			return lexSearch
		}
	}
}

func cobolError(l *lexer) stateFn {
	fmt.Println("Cobol is broken. This tool only works on compilable code.")
	return nil
}
