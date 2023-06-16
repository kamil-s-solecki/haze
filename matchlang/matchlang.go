package matchlang

type Ast interface{}

type NilAst interface{}

var nilast NilAst

type OperatorEnum int

const (
	EqualsOperator OperatorEnum = iota
	NotEqualsOperator
)

type LogicalOperatorEnum int

const (
	AndOperator LogicalOperatorEnum = iota
	OrOperator
)

type IdentifierEnum int

const (
	CodeIdentifier IdentifierEnum = iota
	SizeIdentifier
	TextIdentifier
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

type LogicalExpression struct {
	Operator    LogicalOperatorEnum
	Left, Right Ast
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

func lexTokenToIdentifier(token LexToken) Identifier {
	var idtype IdentifierEnum
	switch token.Type {
	case CodeToken:
		idtype = CodeIdentifier
	case SizeToken:
		idtype = SizeIdentifier
	case TextToken:
		idtype = TextIdentifier
	}
	return Identifier{idtype}
}

func lexTokenToLogicalOperator(token LexToken) LogicalOperatorEnum {
	switch token.Type {
	case AndToken:
		return AndOperator
	case OrToken:
		return OrOperator
	}
	return -1
}

type ParserState int

const (
	ParserConsumingState ParserState = iota
	ParserConsumedLeftState
	ParserConsumedOperatorState
	ParserConsumedRightState
	ParserDoneState
)

type Parser struct {
	tokens           []LexToken
	pos              int
	state            ParserState
	currentLogicalOp LogicalOperatorEnum
	ast              Ast
}

func (p *Parser) consume() bool {
	if p.state == ParserDoneState {
		return false
	}

	switch p.state {
	case ParserConsumingState:
		p.state = ParserConsumedLeftState
	case ParserConsumedLeftState:
		p.state = ParserConsumedOperatorState
	case ParserConsumedOperatorState:
		p.state = ParserConsumedRightState
	case ParserConsumedRightState:
		if p.pos < len(p.tokens)-1 {
			p.currentLogicalOp = lexTokenToLogicalOperator(p.tokens[p.pos])
			p.state = ParserConsumingState
		} else {
			p.state = ParserDoneState
		}
		p.updateAst()
	}
	p.pos++
	return true
}

func (p *Parser) updateAst() {
	ast := Comparison{
		Left:     lexTokenToIdentifier(p.tokens[p.pos-3]),
		Operator: lexTokenToOperator(p.tokens[p.pos-2]),
		Right:    Literal{p.tokens[p.pos-1].Value},
	}
	if p.ast == nilast {
		p.ast = ast
	} else {
		p.ast = LogicalExpression{
			Left:     p.ast,
			Operator: p.currentLogicalOp,
			Right:    ast,
		}
	}
}

func Parse(s string) Ast {
	parser := Parser{tokens: lex(s), pos: 0, state: ParserConsumingState, ast: nilast}
	for parser.consume() {
	}
	return parser.ast
}
