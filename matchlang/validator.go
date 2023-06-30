package matchlang

import (
	"fmt"
	"strings"
)

type ValidatorState int

const (
	ValidatorExpectComparisonState ValidatorState = iota
	ValidatorAfterComparisonState
)

type Validator struct {
	tokens []LexToken
	pos    int
	state  ValidatorState
}

func Validate(expr string) (bool, error) {
	if strings.TrimSpace(expr) == "" {
		return false, fmt.Errorf("The expression cannot be empty!")
	}
	tokens := lex(expr)
	validator := Validator{tokens: tokens}
	err := validator.validate()
	return err == nil, err
}

func (v *Validator) validate() error {
	hasMore, err := v.next()
	for ; hasMore; hasMore, err = v.next() {
	}
	return err
}

func (v *Validator) next() (bool, error) {
	if v.pos >= len(v.tokens) {
		switch v.state {
		case ValidatorExpectComparisonState:
			return false, fmt.Errorf("Expected a comparison after '%v'!", v.tokens[v.pos-1].Value)
		default:
			return false, nil
		}
	}

	if v.isCurrentTokenBracket() {
		v.pos += 1
		return true, nil
	}

	switch v.state {
	case ValidatorExpectComparisonState:
		if err := validateComparison(v.tokens, v.pos); err != nil {
			return false, err
		}
		v.pos += 3
		v.state = ValidatorAfterComparisonState
	case ValidatorAfterComparisonState:
		if err := validateOperator(v.tokens[v.pos]); err != nil {
			return false, err
		}
		v.pos += 1
		v.state = ValidatorExpectComparisonState
	}

	return true, nil
}

func (v *Validator) isCurrentTokenBracket() bool {
	switch v.tokens[v.pos].Type {
	case OpenBracketToken, CloseBracketToken:
		return true
	default:
		return false
	}
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
	if len(tokens) <= idx+1 {
		return fmt.Errorf("Expected an operator after '%v'!", idt.Value)
	}

	op := tokens[idx+1]
	if !isOperator(op) {
		return fmt.Errorf("%v is not a valid operator!", op.Value)
	}
	if len(tokens) <= idx+2 {
		return fmt.Errorf("Expected a literal after '%v'!", op.Value)
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
