package parser

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/walesey/dicelang/dice"

	"github.com/walesey/dicelang/histogram"
	"github.com/walesey/dicelang/token"
)

const ROUND_DECIMALS = 4

type Parser struct {
	lexer Lexer
}

func NewParser(src io.ReadSeeker) Parser {
	return Parser{
		lexer: NewLexer(src),
	}
}

func (parser Parser) Execute() (result interface{}, err error) {
	var operation string
	if _, operation, err = parser.parseToken(token.IDENTIFIER); err != nil {
		return
	}

	var hist histogram.Histogram
	if hist, _, _, err = parser.parseStatement(); err != nil {
		return
	}

	switch operation {
	case "resolve":
		result = hist.Resolve()
	case "hist":
		result = histogram.FormatHistogram(histogram.RoundHistogram(hist, ROUND_DECIMALS))
	case "mean":
		h := hist.Hist()
		var prob float64
		for k, v := range h {
			prob += float64(k) * v
		}
		result = prob
	default:
		err = fmt.Errorf("Invalid operation '%v'", operation)
		return
	}
	return
}

func (parser Parser) parseToken(expect token.Token) (tkn token.Token, literal string, err error) {
	if tkn, literal, err = parser.lexer.Scan(); tkn != expect {
		err = fmt.Errorf("unexpected token: '%v', expected '%v'", tkn.String(), expect.String())
	}
	return
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
		var e error
		if tkn, literal, e = parser.parseToken(token.PERIOD); e != nil {
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

	hist, tkn, literal, err = parser.parseDiceOperator(size, count, false)
	return
}

func (parser Parser) parseDiceOperator(size, count int, invert bool) (hist histogram.Histogram, tkn token.Token, literal string, err error) {
	if tkn, literal, err = parser.parseToken(token.PERIOD); err != nil {
		return
	}

	if tkn, literal, err = parser.parseToken(token.IDENTIFIER); err != nil {
		return
	}

	d := dice.Dice{Size: size}
	if literal == "not" {
		if hist, tkn, literal, err = parser.parseDiceOperator(size, count, !invert); err != nil {
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
	if hist, tkn, literal, err = parser.parseConst(); hist != nil || err != nil {
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

func (parser Parser) parseConst() (hist histogram.Histogram, tkn token.Token, literal string, err error) {
	if tkn, literal, err = parser.lexer.Scan(); err != nil {
		return
	}
	if tkn != token.IDENTIFIER {
		return
	}

	re := regexp.MustCompile(`^([0-9]+)$`)
	matches := re.FindStringSubmatch(literal)
	if len(matches) < 2 {
		return
	}

	var value int
	if value, err = strconv.Atoi(matches[1]); err != nil {
		return
	}

	hist = histogram.Fixed(map[int]float64{value: 1.0})
	return
}
