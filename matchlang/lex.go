package matchlang

import "strings"

type TokenType int

const (
	CodeToken TokenType = iota
	EqualsToken
	LiteralToken
)

type LexToken struct {
	Type  TokenType
	Value string
}

func lex(s string) []LexToken {
	result := []LexToken{}
	for _, word := range strings.Split(s, " ") {
		var token LexToken
		switch word {
		case "code":
			token = LexToken{Type: CodeToken}
		case "=":
			token = LexToken{Type: EqualsToken}
		default:
			token = LexToken{Type: LiteralToken, Value: word}
		}
		result = append(result, token)
	}
	return result
}
