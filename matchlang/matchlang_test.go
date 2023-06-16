package matchlang

import (
	"reflect"
	"testing"
)

func literalEquals(got, want Literal) bool {
	return got.Value == want.Value
}

func identifierEquals(got, want Identifier) bool {
	return got.Value == want.Value
}

func comparisonEquals(got, want Comparison) bool {
	return got.Operator == want.Operator && astEquals(got.Left, want.Left) && astEquals(got.Right, want.Right)
}

func logicalExpressionEquals(got, want LogicalExpression) bool {
	return got.Operator == want.Operator && astEquals(got.Left, want.Left) && astEquals(got.Right, want.Right)
}

func astEquals(got, want Ast) bool {
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		return false
	}
	switch got.(type) {
	case Literal:
		return literalEquals(got.(Literal), want.(Literal))
	case Identifier:
		return identifierEquals(got.(Identifier), want.(Identifier))
	case Comparison:
		return comparisonEquals(got.(Comparison), want.(Comparison))
	case LogicalExpression:
		return logicalExpressionEquals(got.(LogicalExpression), want.(LogicalExpression))
	}
	return false
}

func assertAstEquals(t *testing.T, got, want Ast) {
	if !astEquals(got, want) {
		t.Errorf("got %v, wanted %v", got, want)
	}
}

func TestReturnAstForCodeMatch(t *testing.T) {
	var want Ast
	want = Comparison{Operator: EqualsOperator, Left: Identifier{Value: CodeIdentifier}, Right: Literal{Value: "200"}}

	got := Parse("code = 200")

	assertAstEquals(t, got, want)
}

func TestReturnAstForCodeMatchNotEquals(t *testing.T) {
	var want Ast
	want = Comparison{Operator: NotEqualsOperator, Left: Identifier{Value: CodeIdentifier}, Right: Literal{Value: "200"}}

	got := Parse("code != 200")

	assertAstEquals(t, got, want)
}

func TestReturnAstForSizeMatch(t *testing.T) {
	var want Ast
	want = Comparison{Operator: EqualsOperator, Left: Identifier{Value: SizeIdentifier}, Right: Literal{Value: "1500"}}

	got := Parse("size = 1500")

	assertAstEquals(t, got, want)
}

func TestReturnAstForTextMatch(t *testing.T) {
	var want Ast
	want = Comparison{Operator: EqualsOperator, Left: Identifier{Value: TextIdentifier}, Right: Literal{Value: "foo"}}

	got := Parse("text = foo")

	assertAstEquals(t, got, want)
}

func TestReturnAstWithStringLiteral(t *testing.T) {
	var want Ast
	want = Comparison{Operator: EqualsOperator, Left: Identifier{Value: TextIdentifier}, Right: Literal{Value: "foo"}}

	got := Parse("text = 'foo bar'")

	assertAstEquals(t, got, want)
}

func TestReturnAstWithLogicalExpr(t *testing.T) {
	var want Ast
	left := Comparison{Operator: EqualsOperator, Left: Identifier{Value: CodeIdentifier}, Right: Literal{Value: "200"}}
	right := Comparison{Operator: EqualsOperator, Left: Identifier{Value: SizeIdentifier}, Right: Literal{Value: "1500"}}
	want = LogicalExpression{Operator: AndOperator, Left: left, Right: right}

	got := Parse("code = 200 and size = 1500")

	assertAstEquals(t, got, want)
}
