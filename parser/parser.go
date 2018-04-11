package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/walesey/dicelang/dice"

	"github.com/walesey/dicelang/histogram"
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

func (parser Parser) Execute() (result string, err error) {
	var operation string
	if operation, err = parser.parseToken(token.IDENTIFIER); err != nil {
		return
	}

	if _, err = parser.parseToken(token.WHITESPACE); err != nil {
		return
	}

	var hist histogram.Histogram
	if hist, err = parser.parseStatement(); err != nil {
		return
	}

	var data interface{}
	switch operation {
	case "resolve":
		data = hist.Resolve()
	case "hist":
		data = histogram.RoundHistogram(hist.Hist())
	case "mean":
		h := hist.Hist()
		var prob float64
		for k, v := range h {
			prob += float64(k) * v
		}
		data = prob
	default:
		err = fmt.Errorf("Invalid operation '%v'", operation)
		return
	}

	var resultData []byte
	resultData, err = json.Marshal(data)
	result = string(resultData)
	return
}

func (parser Parser) parseToken(expect token.Token) (string, error) {
	tkn, literal, err := parser.lexer.Scan()
	if err != nil {
		return "", err
	} else if tkn != expect {
		err = fmt.Errorf("unexpected token: '%v', expected '%v'", tkn.String(), expect.String())
	}
	return literal, nil
}

func (parser Parser) parseStatement() (histogram.Histogram, error) {
	histograms := []histogram.Histogram{}
	for {
		hist, _, _, err := parser.parseDice()
		if err != nil {
			return nil, err
		}
		if hist == nil {
			break
		}
		histograms = append(histograms, hist)
	}
	return histogram.Multiply(histograms...), nil
}

func (parser Parser) parseDice() (hist histogram.Histogram, tkn token.Token, literal string, err error) {
	if hist, tkn, literal, err = parser.parseAggregate(); hist != nil || err != nil {
		return
	}
	if tkn != token.IDENTIFIER {
		return
	}

	re := regexp.MustCompile(`(|[0-9]+)d([0-9]+)`)
	matches := re.FindStringSubmatch(literal)
	fmt.Println(matches)

	if len(matches) < 3 {
		err = fmt.Errorf("Invalid dice syntax: '%v'", literal)
		return
	}

	size, count := 0, 0
	if matches[1] == "" {
		count = 1
	} else if count, err = strconv.Atoi(matches[1]); err != nil {
		return
	}

	if matches[2] == "" {
		size = 1
	} else if size, err = strconv.Atoi(matches[2]); err != nil {
		return
	}

	hist, err = parser.parseDiceOperator(size, count)
	return
}

func (parser Parser) parseDiceOperator(size, count int) (hist histogram.Histogram, err error) {
	d := dice.Dice{Size: size}

	if _, err = parser.parseToken(token.PERIOD); err != nil {
		return
	}

	var literal string
	if literal, err = parser.parseToken(token.IDENTIFIER); err != nil {
		return
	}

	if literal == "not" {
		if hist, err = parser.parseDiceOperator(size, count); err != nil {
			return
		}
		hist = histogram.Invert(hist)
	} else if literal == "add" {
		hist = dice.MultiDice{Dice: d, Count: count}
	} else {
		re := regexp.MustCompile(`([0-9]+)\+`)
		matches := re.FindStringSubmatch(literal)
		var gte int
		if gte, err = strconv.Atoi(matches[1]); err != nil {
			return
		}
		hist = dice.DicePool{Dice: d, Count: count, GTE: gte}
	}

	return
}

func (parser Parser) parseAggregate() (hist histogram.Histogram, tkn token.Token, literal string, err error) {
	if tkn, literal, err = parser.lexer.Scan(); err != nil {
		return
	}
	if tkn != token.OPEN_BRACKET {
		return
	}

	//TODO

	return
}
