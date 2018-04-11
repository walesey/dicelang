package parser

import (
	"fmt"
	"io"

	"github.com/walesey/dicelang/token"
)

type Parser struct {
	lexer Lexer
}

func NewParser(src io.ReadSeeker) Parser {
	return Parser{
		lexer: NewLexer(src),
	}
}

func (parser Parser) parseToken(expect token.Token) (string, error) {
	tkn, literal, err := parser.lexer.Scan()
	if err != nil {
		return "", err
	} else if tkn != token.IDENTIFIER {
		err = fmt.Errorf("unexpected token: '%v', expected '%v'", tkn.String(), expect.String())
	}
	return literal, nil
}

func (parser Parser) parseStatement() (hist Histogram, err error) {
	var literal string
	if literal, err = parser.parseToken(token.IDENTIFIER); err != nil {
		return
	}

}

func (parser Parser) SimpleExec() (result string, err error) {
	var operation string
	if operation, err = parser.parseToken(token.IDENTIFIER); err != nil {
		return
	}

	if _, err = parser.parseToken(token.WHITESPACE); err != nil {
		return
	}

	var hist Histogram
	if hist, err = parser.parseStatement(); err != nil {
		return
	}

	var data interface{}
	switch operation {
	case "resolve":
		data = hist.Resolve()
	case "hist":
		data = hist.Hist()
	case "mean":
		h := hist.Hist()
		var prob float64
		for k, v := range h {
			prob += float64(k) * v
		}
		data = prob
	default:
		err := fmt.Errorf("Invalid operation '%v'", operation)
	}

	return
}
