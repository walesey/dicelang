package token

// Token is the set of lexical tokens in dicelang
type Token int

const (
	_ Token = iota

	ILLEGAL
	EOF
	WHITESPACE
	COMMENT // //
	KEYWORD

	STRING
	BOOLEAN
	NUMBER
	IDENTIFIER

	OPEN_BRACKET  // [
	CLOSE_BRACKET // ]
	COMMA         // ,
	PERIOD        // .
)

func (tkn Token) String() string {
	switch tkn {
	case WHITESPACE:
		return "WHITESPACE"
	case STRING:
		return "STRING"
	case BOOLEAN:
		return "BOOLEAN"
	case NUMBER:
		return "NUMBER"
	case IDENTIFIER:
		return "IDENTIFIER"
	case OPEN_BRACKET:
		return "["
	case CLOSE_BRACKET:
		return "]"
	case COMMA:
		return ","
	case PERIOD:
		return "."
	default:
		return ""
	}
}
