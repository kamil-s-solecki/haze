package matchlang

type Ast interface{}

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
	tokens              []LexToken
	pos                 int
	state               ParserState
	currentLogicalOp    LogicalOperatorEnum
	isLastExprBracketed bool
	ast                 Ast
}

func (p *Parser) consume() bool {
	if p.state == ParserDoneState {
		return false
	}

	switch p.state {
	case ParserConsumingState:
		if !p.isCurrentTokenBracket() {
			p.state = ParserConsumedLeftState
		}
	case ParserConsumedLeftState:
		p.state = ParserConsumedOperatorState
	case ParserConsumedOperatorState:
		p.state = ParserConsumedRightState
	case ParserConsumedRightState:
		p.updateAst()
		if p.canConsumeMore() {
			p.readTillLogicalOp()
			p.state = ParserConsumingState
		} else {
			p.state = ParserDoneState
		}
	}
	p.pos++
	return true
}

func (p *Parser) updateAst() {
	switch p.ast.(type) {
	case nil:
		p.ast = p.currentLeaf()
	case LogicalExpression:
		if p.ast.(LogicalExpression).Operator == AndOperator || p.isLastExprBracketed {
			p.addNodeAbove()
		} else {
			p.addNodeBelow()
		}
	case Comparison:
		p.addNodeAbove()
	}
	p.isLastExprBracketed = false
}

func (p *Parser) addNodeAbove() {
	p.ast = LogicalExpression{
		Left:     p.ast,
		Operator: p.currentLogicalOp,
		Right:    p.currentLeaf(),
	}
}

func (p *Parser) addNodeBelow() {
	oldAst := p.ast.(LogicalExpression)
	right := LogicalExpression{
		Left:     oldAst.Right,
		Operator: p.currentLogicalOp,
		Right:    p.currentLeaf(),
	}
	p.ast = LogicalExpression{
		Left:     oldAst.Left,
		Operator: oldAst.Operator,
		Right:    right,
	}
}

func (p *Parser) currentLeaf() Comparison {
	return Comparison{
		Left:     lexTokenToIdentifier(p.tokens[p.pos-3]),
		Operator: lexTokenToOperator(p.tokens[p.pos-2]),
		Right:    Literal{p.tokens[p.pos-1].Value},
	}
}

func (p *Parser) currentToken() LexToken {
	return p.tokens[p.pos]
}

func (p *Parser) isCurrentTokenBracket() bool {
	switch p.currentToken().Type {
	case OpenBracketToken, CloseBracketToken:
		return true
	default:
		return false
	}
}

func (p *Parser) readTillLogicalOp() {
	if p.isCurrentTokenBracket() {
		p.isLastExprBracketed = true
		p.pos++
	}
	p.currentLogicalOp = lexTokenToLogicalOperator(p.currentToken())
}

func (p *Parser) canConsumeMore() bool {
	return p.pos < len(p.tokens)-1
}

func Parse(s string) Ast {
	parser := Parser{tokens: lex(s), state: ParserConsumingState}
	for parser.consume() {
	}
	return parser.ast
}
