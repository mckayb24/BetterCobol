package BetterCobol

import "testing"

func TestLexer(t *testing.T) {
	tests := []struct {
		input      string
		bufferSize int
		expected   []item
		reason     string
	}{
		{
			`these are random    
			words that 
			don't matter at all
			EXEC SQL INCLUDE test END-EXEC.
			these are more random
			words`,
			0,
			[]item{
				{itemCopyBook, "test"},
				{itemEOF, ""},
			},
			"Copy book statement surrounded by garbage",
		},
		{
			`these are random
       DATA DIVISION.
        content goes here`,
			0,
			[]item{
				{itemDataDivision, "\n        content goes here"},
				{itemEOF, ""},
			},
			"Data division end by eof",
		},
		{
			`these are random
       DATA DIVISION.
        content goes here
       PROCEDURE DIVISION.
        procedure content here`,
			0,
			[]item{
				{itemDataDivision, "\n        content goes here"},
				{itemProcedureDivision, "\n        procedure content here"},
				{itemEOF, ""},
			},
			"Data division end by procedure division",
		},
		{
			`these are random
       DATA DIVISION.
        content goes here
       PROCEDURE DIVISION USING stuff and more  .
        procedure content here`,
			0,
			[]item{
				{itemDataDivision, "\n        content goes here"},
				{itemProcedureDivision, "\n        procedure content here"},
				{itemEOF, ""},
			},
			"Data division end by procedure division with using",
		},
	}
	for _, test := range tests {
		_, ch := lex("test", test.input, test.bufferSize)
		for _, itm := range test.expected {
			actual := <-ch
			if itm != actual {
				t.Fatalf("expected: %v, actual: %v, reason: %s", itm, actual, test.reason)
			}
		}
	}
}
