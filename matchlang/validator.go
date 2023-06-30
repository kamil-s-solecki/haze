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

	if err := validateComparison(tokens, 0); err != nil {
		return false, err
	}

	if len(tokens) > 3 {
		if err := validateOperator(tokens[3]); err != nil {
			return false, err
		}

		if err := validateComparison(tokens, 4); err != nil {
			return false, err
		}
	}

	return true, nil
}

func validateOperator(token LexToken) error {
	switch token.Type {
	case OrToken, AndToken:
		return nil
	default:
		return fmt.Errorf("%v is not a valid logical operator!", token.Value)
	}
}

func validateComparison(tokens []LexToken, idx int) error {
	idt := tokens[idx+0]
	if !isIdentifier(idt) {
		return fmt.Errorf("%v is not a valid identifier!", idt.Value)
	}

	op := tokens[idx+1]
	if !isOperator(op) {
		return fmt.Errorf("%v is not a valid operator!", op.Value)
	}

	lit := tokens[idx+2]
	if !isLiteral(lit) {
		return fmt.Errorf("%v is not a valid literal!", lit.Value)
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
