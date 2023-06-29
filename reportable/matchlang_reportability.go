package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/matchlang"
)

type Checker func(http.Response) bool

func Compile(expr string) Checker {
	ast := matchlang.Parse(expr)
	matcher := comparisonToChecker(ast.(matchlang.Comparison))
	return func(r http.Response) bool {
		return matcher(r)
	}
}

func comparisonToChecker(comp matchlang.Comparison) Checker {
	val := comp.Right.(matchlang.Literal).Value
	id := comp.Left.(matchlang.Identifier).Value
	if comp.Operator == matchlang.EqualsOperator {
		switch id {
		case matchlang.CodeIdentifier:
			return Checker(MatchCodes(val))
		case matchlang.SizeIdentifier:
			return Checker(MatchLengths(val))
		case matchlang.TextIdentifier:
			return Checker(MatchString(val))
		}
	} else {
		switch id {
		case matchlang.CodeIdentifier:
			return Checker(FilterCodes(val))
		case matchlang.SizeIdentifier:
			return Checker(FilterLengths(val))
		case matchlang.TextIdentifier:
			return Checker(FilterString(val))
		}
	}
	return nil
}
