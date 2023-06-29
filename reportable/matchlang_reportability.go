package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/matchlang"
)

func Compile(expr string) func(http.Response) bool {
	ast := matchlang.Parse(expr)
	matcher := comparisonToMatcher(ast.(matchlang.Comparison))
	return func(r http.Response) bool {
		return matcher(r)
	}
}

func comparisonToMatcher(comp matchlang.Comparison) Matcher {
	return MatchCodes(comp.Right.(matchlang.Literal).Value)
}
