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
	if len(tokens) == 0 {
		return false, fmt.Errorf("The expression cannot be empty!")
	}
	switch tokens[0].Type {
	case CodeToken, SizeToken, TextToken:
		return true, nil
	default:
		return false, fmt.Errorf("%v is not a valid identifier!", tokens[0].Value)
	}
}
