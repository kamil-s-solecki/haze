package matchlang

type Ast interface{}

type OperatorEnum int

const (
	EqualsOperator OperatorEnum = iota
	NotEqualsOperator
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

func lexTokenToOperator(token LexToken) OperatorEnum {
	switch token.Type {
	case EqualsToken:
		return EqualsOperator
	case NotEqualsToken:
		return NotEqualsOperator
	}
	return -1
}

func Parse(s string) Ast {
	tokens := lex(s)
	return Comparison{
		Left:     Identifier{CodeIdentifier},
		Operator: lexTokenToOperator(tokens[1]),
		Right:    Literal{tokens[2].Value},
	}
}
