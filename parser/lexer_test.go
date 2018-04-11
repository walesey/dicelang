package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walesey/dicelang/token"
)

func Test_Lexer_Read(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte("hist 2d3.4+")))
	c, e := lexer.read()
	assert.Nil(t, e)
	assert.EqualValues(t, 'h', c)
	c, e = lexer.read()
	assert.Nil(t, e)
	assert.EqualValues(t, 'i', c)
}

func Test_Lexer_Peek(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte("hist 2d3.4+")))
	c, e := lexer.read()
	assert.Nil(t, e)
	assert.EqualValues(t, 'h', c)
	c, e = lexer.peek()
	assert.Nil(t, e)
	assert.EqualValues(t, 'i', c)
	c, e = lexer.read()
	assert.Nil(t, e)
	assert.EqualValues(t, 'i', c)
}

func Test_Lexer_ParseIdentifier(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte("hist 2d3.4+")))
	i, e := lexer.parseIdentifier()
	assert.Nil(t, e)
	assert.EqualValues(t, "hist", i)

	c, e := lexer.read()
	assert.Nil(t, e)
	assert.EqualValues(t, ' ', c)

	i, e = lexer.parseIdentifier()
	assert.Nil(t, e)
	assert.EqualValues(t, "2d3", i)

	c, e = lexer.read()
	assert.Nil(t, e)
	assert.EqualValues(t, '.', c)

	i, e = lexer.parseIdentifier()
	assert.Nil(t, e)
	assert.EqualValues(t, "4+", i)
}

func Test_Lexer_Scan(t *testing.T) {
	lexer := NewLexer(bytes.NewReader([]byte("hist 2d3.4+")))
	tkn, l, e := lexer.Scan()
	assert.Nil(t, e)
	assert.EqualValues(t, token.IDENTIFIER, tkn)
	assert.EqualValues(t, "hist", l)

	tkn, l, e = lexer.Scan()
	assert.Nil(t, e)
	assert.EqualValues(t, token.WHITESPACE, tkn)
	assert.EqualValues(t, " ", l)

	tkn, l, e = lexer.Scan()
	assert.Nil(t, e)
	assert.EqualValues(t, token.IDENTIFIER, tkn)
	assert.EqualValues(t, "2d3", l)

	tkn, l, e = lexer.Scan()
	assert.Nil(t, e)
	assert.EqualValues(t, token.PERIOD, tkn)
	assert.EqualValues(t, "", l)

	tkn, l, e = lexer.Scan()
	assert.Nil(t, e)
	assert.EqualValues(t, token.IDENTIFIER, tkn)
	assert.EqualValues(t, "4+", l)
}
