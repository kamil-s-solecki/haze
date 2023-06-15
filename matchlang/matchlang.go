package matchlang

import (
	"strings"
)

type Ast interface{}

type OperatorEnum int

const (
	Equals OperatorEnum = iota
)

type IdentifierEnum int

const (
	CodeIdentifier IdentifierEnum = iota
)

type Comparison struct {
	Operator    OperatorEnum
	Left, Right Ast
}

type Identifier struct {
	Value IdentifierEnum
}

type Literal struct {
	Value string
}

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

func Parse(s string) Ast {
	tokens := lex(s)
	return Comparison{
		Operator: Equals,
		Left:     Identifier{CodeIdentifier},
		Right:    Literal{tokens[2].Value},
	}
}
