package parser_test

import (
	"testing"

	"github.com/romshark/gapi/compiler/parser"
	"github.com/stretchr/testify/require"
)

func isTestFile(t *testing.T, fl *parser.File) {
	require.NotNil(t, fl)
	require.Equal(t, "test.schema", fl.Name)
	require.Equal(t, "/tests/", fl.Path)
}

func compareCursor(t *testing.T, expected, actual parser.Cursor) {
	require.Equal(t, expected.Index, actual.Index)
	require.Equal(t, expected.Line, actual.Line)
	require.Equal(t, expected.Column, actual.Column)
	isTestFile(t, actual.File)
}

func TestLexer(t *testing.T) {
	test := func(
		t *testing.T,
		source string,
		expectedSrc string,
		expectedFragID parser.FragID,
		expectedBegin parser.Cursor,
		expectedEnd parser.Cursor,
	) {
		lexer := parser.NewLexer(src(source))
		require.NotNil(t, lexer)
		tk, err := lexer.Next()
		require.NotNil(t, tk)
		require.NoError(t, err)
		require.Equal(t, expectedFragID, tk.FragID())
		require.Equal(t, expectedSrc, tk.Src())
		compareCursor(t, expectedBegin, tk.Begin())
		compareCursor(t, expectedEnd, tk.End())
	}

	t.Run("singleSpace", func(t *testing.T) {
		test(
			t,
			" f",
			" ",
			parser.FragTkSpace,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 1, Line: 1, Column: 2},
		)
	})
	t.Run("spaces", func(t *testing.T) {
		test(
			t,
			"    f",
			"    ",
			parser.FragTkSpace,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 4, Line: 1, Column: 5},
		)
	})
	t.Run("spacesWithTab", func(t *testing.T) {
		test(
			t,
			"  \t  f ",
			"  \t  ",
			parser.FragTkSpace,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 5, Line: 1, Column: 6},
		)
	})
	t.Run("spaceWithTabsAndLineBreak", func(t *testing.T) {
		test(
			t,
			"  \t \n\t  f ",
			"  \t \n\t  ",
			parser.FragTkSpace,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 8, Line: 2, Column: 4},
		)
	})
	t.Run("parser.FragTkLatinAlphanum", func(t *testing.T) {
		test(
			t,
			"word",
			"word",
			parser.FragTkLatinAlphanum,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 4, Line: 1, Column: 5},
		)
	})
	t.Run("SymOpt", func(t *testing.T) {
		test(
			t,
			"?",
			"?",
			parser.FragTkSymOpt,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 1, Line: 1, Column: 2},
		)
	})
}

// TestLexerScanSequence tests scanning a sequence of tokens
func TestLexerScanSequence(t *testing.T) {
	tkz := parser.NewLexer(src("  \t \n  word\r\n 345 x\t\t\t"))
	require.NotNil(t, tkz)

	type Token struct {
		id    parser.FragID
		src   string
		begin parser.Cursor
		end   parser.Cursor
	}
	expected := []Token{
		Token{
			parser.FragTkSpace,
			"  \t \n  ",
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 7, Line: 2, Column: 3},
		},
		Token{
			parser.FragTkLatinAlphanum,
			"word",
			parser.Cursor{Index: 7, Line: 2, Column: 3},
			parser.Cursor{Index: 11, Line: 2, Column: 7},
		},
		Token{
			parser.FragTkSpace,
			"\r\n ",
			parser.Cursor{Index: 11, Line: 2, Column: 7},
			parser.Cursor{Index: 14, Line: 3, Column: 2},
		},
		Token{
			parser.FragTkLatinAlphanum,
			"345",
			parser.Cursor{Index: 14, Line: 3, Column: 2},
			parser.Cursor{Index: 17, Line: 3, Column: 5},
		},
		Token{
			parser.FragTkSpace,
			" ",
			parser.Cursor{Index: 17, Line: 3, Column: 5},
			parser.Cursor{Index: 18, Line: 3, Column: 6},
		},
		Token{
			parser.FragTkLatinAlphanum,
			"x",
			parser.Cursor{Index: 18, Line: 3, Column: 6},
			parser.Cursor{Index: 19, Line: 3, Column: 7},
		},
		Token{
			parser.FragTkSpace,
			"\t\t\t",
			parser.Cursor{Index: 19, Line: 3, Column: 7},
			parser.Cursor{Index: 22, Line: 3, Column: 10},
		},
	}

	for _, expected := range expected {
		tk, err := tkz.Next()
		require.NoError(t, err)
		require.NotNil(t, tk)
		require.Equal(t, expected.src, tk.Src())
		compareCursor(t, expected.begin, tk.Begin())
		compareCursor(t, expected.end, tk.End())
		require.Equal(t, expected.id, tk.FragID())
	}

	tk, err := tkz.Next()
	require.NoError(t, err)
	require.Nil(t, tk)
}

// TestLexerSyntaxErr tests lexer syntax errors
func TestLexerSyntaxErr(t *testing.T) {
	test := func(
		t *testing.T,
		source string,
		expectedAt parser.Cursor,
	) {
		lexer := parser.NewLexer(src(source))
		require.NotNil(t, lexer)

		tk, err := lexer.Next()
		require.Error(t, err)
		require.Nil(t, tk)

		require.Equal(t, parser.ErrSyntax, err.Code())
		compareCursor(t, expectedAt, err.At())
	}

	t.Run("UnclosedList", func(t *testing.T) {
		test(
			t,
			"[",
			parser.Cursor{Index: 0, Line: 1, Column: 1},
		)
	})
	t.Run("ReturnCarriageWithoutLineFeed", func(t *testing.T) {
		test(
			t,
			"\r",
			parser.Cursor{Index: 0, Line: 1, Column: 1},
		)
	})
}

// TestLexerNextExpect tests lexer syntax errors
func TestLexerNextExpect(t *testing.T) {
	test := func(
		t *testing.T,
		source,
		errMsg,
		expectedErrMsg string,
		expectedFragID parser.FragID,
		expectedBegin parser.Cursor,
	) {
		lexer := parser.NewLexer(src(source))
		require.NotNil(t, lexer)
		tk, err := lexer.NextExpect(expectedFragID, errMsg)
		require.Error(t, err)
		require.Equal(t, expectedErrMsg, err.Message())
		if tk != nil {
			compareCursor(t, tk.Begin(), err.At())
			require.NotEqual(t, expectedFragID, tk.FragID())
		}
	}

	t.Run("expectLatinAlphanum", func(t *testing.T) {
		test(
			t,
			".",
			"expected alpha-num",
			"expected alpha-num, got: '.'",
			parser.FragTkLatinAlphanum,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
		)
	})
	t.Run("unexpectedEOF", func(t *testing.T) {
		test(
			t,
			"",
			"expected alpha-num",
			"expected alpha-num, reached end of file",
			parser.FragTkLatinAlphanum,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
		)
	})
}

// TestLexerNextExpectSkip tests lexer syntax errors
func TestLexerNextExpectSkip(t *testing.T) {
	test := func(
		t *testing.T,
		source string,
		skip parser.Skip,
		expectedFragID parser.FragID,
		expectedSrc string,
		expectedBegin parser.Cursor,
		expectedEnd parser.Cursor,
	) {
		lexer := parser.NewLexer(src(source))
		require.NotNil(t, lexer)
		tk, err := lexer.NextExpectSkip(
			expectedFragID,
			skip,
			"expected %s",
			expectedSrc,
		)
		require.NoError(t, err)
		require.NotNil(t, tk)
		require.Equal(t, expectedFragID, tk.FragID())
		require.Equal(t, expectedSrc, tk.Src())
		compareCursor(t, expectedBegin, tk.Begin())
		compareCursor(t, expectedEnd, tk.End())
	}

	t.Run("skipSpace", func(t *testing.T) {
		test(
			t,
			" \n alpha",
			parser.Skip{parser.FragTkSpace},
			parser.FragTkLatinAlphanum,
			"alpha",
			parser.Cursor{Index: 3, Line: 2, Column: 2},
			parser.Cursor{Index: 8, Line: 2, Column: 7},
		)
	})
	t.Run("noSkip", func(t *testing.T) {
		test(
			t,
			"alpha",
			parser.Skip{parser.FragTkSpace},
			parser.FragTkLatinAlphanum,
			"alpha",
			parser.Cursor{Index: 0, Line: 1, Column: 1},
			parser.Cursor{Index: 5, Line: 1, Column: 6},
		)
	})
	t.Run("skipMultiple", func(t *testing.T) {
		test(
			t,
			"  []alpha, ?",
			parser.Skip{
				parser.FragTkSpace,
				parser.FragTkLatinAlphanum,
				parser.FragTkSymList,
				parser.FragTkSymSep,
			},
			parser.FragTkSymOpt,
			"?",
			parser.Cursor{Index: 11, Line: 1, Column: 12},
			parser.Cursor{Index: 12, Line: 1, Column: 13},
		)
	})
}

// TestLexerNextExpectSkipErr tests lexer syntax errors
func TestLexerNextExpectSkipErr(t *testing.T) {
	test := func(
		t *testing.T,
		source string,
		skip parser.Skip,
		errMsg,
		expectedErrMsg string,
		expectedFragID parser.FragID,
		expectedBegin parser.Cursor,
	) {
		lexer := parser.NewLexer(src(source))
		require.NotNil(t, lexer)
		tk, err := lexer.NextExpectSkip(expectedFragID, skip, errMsg)
		require.Error(t, err)
		require.Equal(t, expectedErrMsg, err.Message())
		if tk != nil {
			compareCursor(t, tk.Begin(), err.At())
			require.NotEqual(t, expectedFragID, tk.FragID())
		}
	}
	t.Run("unexpectedEOF", func(t *testing.T) {
		test(
			t,
			"",
			parser.Skip{parser.FragTkSpace},
			"expected alpha-num",
			"expected alpha-num, reached end of file",
			parser.FragTkLatinAlphanum,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
		)
	})
	t.Run("unexpectedEOF", func(t *testing.T) {
		test(
			t,
			" alpha",
			parser.Skip{
				parser.FragTkSpace,
				parser.FragTkLatinAlphanum,
			},
			"expected alpha-num",
			"expected alpha-num, reached end of file",
			parser.FragTkSymSep,
			parser.Cursor{Index: 0, Line: 1, Column: 1},
		)
	})
}
