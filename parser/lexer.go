package parser

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/walesey/dicelang/token"

	"io"
)

type EofError struct{ error }
type UnexpectedTokenError struct{ error }
type InvalidIdentifierStartError struct{ error }

type RuneFunc func(chr rune) bool

type Lexer struct {
	src        io.ReadSeeker
	readBuffer []byte
}

func NewLexer(src io.ReadSeeker) Lexer {
	return Lexer{
		src:        src,
		readBuffer: make([]byte, 1),
	}
}

func (lexer Lexer) read() (chr rune, err error) {
	var n int
	if n, err = lexer.src.Read(lexer.readBuffer); n == 0 {
		return ' ', EofError{fmt.Errorf("EOF")}
	} else if err != nil {
		return
	}

	chr = rune(lexer.readBuffer[0])
	return
}

func (lexer Lexer) peek() (chr rune, err error) {
	if chr, err = lexer.read(); err != nil {
		return
	}

	_, err = lexer.src.Seek(-1, io.SeekCurrent)
	return
}

func (lexer Lexer) Scan() (tkn token.Token, literal string, err error) {
	if literal, err = lexer.parseIdentifier(); err == nil {
		tkn = token.IDENTIFIER
		return
	}

	var chr rune
	if chr, err = lexer.read(); err != nil {
		if _, ok := err.(EofError); ok {
			tkn, err = token.EOF, nil
		}
		return
	}

	switch chr {
	case '/':
		if chr, err = lexer.read(); err != nil {
			return
		} else if chr == '/' {
			lexer.read()
			tkn = token.COMMENT
		}
	case ' ':
		if _, err = lexer.src.Seek(-1, io.SeekCurrent); err != nil {
			return
		}
		if literal, err = lexer.parseWhiteSpace(); err != nil {
			return
		}
		tkn = token.WHITESPACE
	case '[':
		tkn = token.OPEN_BRACKET
	case ']':
		tkn = token.CLOSE_BRACKET
	case '.':
		tkn = token.PERIOD
	case ',':
		tkn = token.COMMA
	case '\n', '\r':
		tkn = token.NEWLINE
	default:
		err = UnexpectedTokenError{fmt.Errorf(fmt.Sprintf("Unexpected Token '%c'", chr))}
		tkn = token.ILLEGAL
	}

	return
}

func (lexer Lexer) parseGeneric(isPart RuneFunc) (literal string, err error) {
	var chr rune
	if chr, err = lexer.peek(); err != nil {
		return
	}

	var offset, strSize int64 = 0, -1
	for ; isPart(chr); offset++ {
		if chr, err = lexer.read(); err != nil {
			if _, ok := err.(EofError); ok {
				strSize++
				break
			}
			return
		}
	}

	// Seek back and read the string
	if _, err = lexer.src.Seek(-offset, io.SeekCurrent); err != nil {
		return
	}
	strSize += offset
	b := make([]byte, strSize)
	if _, err = lexer.src.Read(b); err != nil {
		return
	}
	literal = string(b)
	return
}

func (lexer Lexer) parseWhiteSpace() (literal string, err error) {
	return lexer.parseGeneric(isWhiteSpace)
}

func (lexer Lexer) parseIdentifier() (literal string, err error) {
	var chr rune
	if chr, err = lexer.peek(); err != nil {
		return
	} else if !isIdentifierStart(chr) {
		err = InvalidIdentifierStartError{fmt.Errorf("Invalid character used at the start of identifier")}
		return
	}

	return lexer.parseGeneric(isIdentifierPart)
}

func isWhiteSpace(chr rune) bool {
	return chr == ' ' || chr == '\t'
}

func isIdentifierStart(chr rune) bool {
	return isIdentifierPart(chr)
}

func isIdentifierPart(chr rune) bool {
	return chr == '+' || chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		'0' <= chr && chr <= '9' ||
		chr >= utf8.RuneSelf && (unicode.IsLetter(chr) || unicode.IsDigit(chr))
}
