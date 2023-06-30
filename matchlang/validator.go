package matchlang

import (
	"fmt"
	"strings"
)

func Validate(expr string) (bool, error) {
	if strings.TrimSpace(expr) == "" {
		return false, fmt.Errorf("The expression cannot be empty!")
	}

	tokens := lex(expr)

	if err := validateComparison(tokens); err != nil {
		return false, err
	}

	return true, nil
}

func validateComparison(tokens []LexToken) error {
	if !isIdentifier(tokens[0]) {
		return fmt.Errorf("%v is not a valid identifier!", tokens[0].Value)
	}

	if !isOperator(tokens[1]) {
		return fmt.Errorf("%v is not a valid operator!", tokens[1].Value)
	}

	if !isLiteral(tokens[2]) {
		return fmt.Errorf("%v is not a valid literal!", tokens[2].Value)
	}

	return nil
}

func isIdentifier(token LexToken) bool {
	switch token.Type {
	case CodeToken, SizeToken, TextToken:
		return true
	default:
		return false
	}
}

func isOperator(token LexToken) bool {
	switch token.Type {
	case EqualsToken, NotEqualsToken:
		return true
	default:
		return false
	}
}

func isLiteral(token LexToken) bool {
	switch token.Type {
	case LiteralToken:
		return true
	default:
		return false
	}
}
