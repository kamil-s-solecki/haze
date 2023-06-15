package matchlang

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

func Parse(s string) Ast {
	tokens := lex(s)
	return Comparison{
		Operator: Equals,
		Left:     Identifier{CodeIdentifier},
		Right:    Literal{tokens[2].Value},
	}
}
