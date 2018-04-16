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

	var hist histogram.Histogram
	if hist, _, _, err = parser.parseStatement(); err != nil {
		return
	}

	var data interface{}
	switch operation {
	case "resolve":
		data = hist.Resolve()
	case "hist":
		data = histogram.FormatHistogram(histogram.RoundHistogram(hist, 3))
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
		return "", fmt.Errorf("unexpected token: '%v', expected '%v'", tkn.String(), expect.String())
	}
	return literal, nil
}

func (parser Parser) parseStatement() (hist histogram.Histogram, tkn token.Token, literal string, err error) {
	histograms := []histogram.Histogram{}
	for {
		hist, tkn, literal, err = parser.parseDice()
		if err != nil {
			return
		}
		if hist == nil {
			err = fmt.Errorf("Invalid syntax")
			return
		}
		histograms = append(histograms, histogram.RoundHistogram(hist, 9))
		if _, e := parser.parseToken(token.PERIOD); e != nil {
			break
		}
	}
	hist = histogram.Multiply(histograms...)
	return
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

	hist, err = parser.parseDiceOperator(size, count, false)
	return
}

func (parser Parser) parseDiceOperator(size, count int, invert bool) (hist histogram.Histogram, err error) {
	if _, err = parser.parseToken(token.PERIOD); err != nil {
		return
	}

	var tkn token.Token
	var literal string
	if tkn, literal, err = parser.lexer.Scan(); err != nil {
		return
	}

	switch tkn {
	// case token.OPEN_BRACKET:
	// 	if hist, err = parser.parseDiceOperatorAggregate(size, count, invert); err != nil {
	// 		return
	// 	}
	case token.IDENTIFIER:
		if hist, err = parser.parseDiceOperatorLiteral(size, count, invert, literal); err != nil {
			return
		}
	default:
		err = fmt.Errorf("Unexpected token '%v', expected dice operator", tkn.String())
	}

	return
}

// func (parser Parser) parseDiceOperatorAggregate(size, count int, invert bool) (hist histogram.Histogram, err error) {
// 	histograms := []histogram.Histogram{}

// 	var tkn token.Token
// 	var literal string
// 	for {
// 		if literal, err = parser.parseToken(token.IDENTIFIER); err != nil {
// 			return
// 		}
// 		if hist, err = parser.parseDiceOperatorLiteral(size, count, invert, literal); err != nil {
// 			return
// 		}
// 		if hist, tkn, literal, err = parser.parseStatement(); err != nil {
// 			return
// 		}
// 		histograms = append(histograms, hist)
// 		if tkn != token.COMMA {
// 			break
// 		}
// 	}

// 	hist = histogram.Aggregate(histograms...)
// 	if tkn != token.CLOSE_BRACKET {
// 		err = fmt.Errorf("Unexpected token '%v', expected close bracket ']'", tkn.String())
// 	}
// 	return
// }

func (parser Parser) parseDiceOperatorLiteral(size, count int, invert bool, literal string) (hist histogram.Histogram, err error) {
	d := dice.Dice{Size: size}
	if literal == "not" {
		if hist, err = parser.parseDiceOperator(size, count, !invert); err != nil {
			return
		}
	} else if literal == "add" {
		hist = dice.MultiDice{Dice: d, Count: count}
	} else { // gte: eg. 4+, 6+, ...
		re := regexp.MustCompile(`([0-9]+)\+`)
		matches := re.FindStringSubmatch(literal)
		if len(matches) < 2 {
			err = fmt.Errorf("Invalid dice operator syntax: '%v'", literal)
			return
		}
		var gte int
		if gte, err = strconv.Atoi(matches[1]); err != nil {
			return
		}
		if invert {
			hist = dice.DicePool{Dice: d, Count: count, GTE: size - gte + 2}
		} else {
			hist = dice.DicePool{Dice: d, Count: count, GTE: gte}
		}
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

	histograms := []histogram.Histogram{}
	for {
		if hist, tkn, literal, err = parser.parseStatement(); err != nil {
			return
		}
		histograms = append(histograms, hist)
		if tkn != token.COMMA {
			break
		}
	}

	hist = histogram.Aggregate(histograms...)
	if tkn != token.CLOSE_BRACKET {
		err = fmt.Errorf("Unexpected token '%v', expected close bracket ']'", tkn.String())
	}
	return
}
