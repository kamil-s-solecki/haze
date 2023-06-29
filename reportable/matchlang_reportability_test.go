package reportable

import (
	"github.com/kamil-s-solecki/haze/http"
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestCompiledShouldNotReport200(t *testing.T) {
	res := http.Response{Code: 200}
	checker := Compile("code = 500")

	got := checker(res)

	testutils.AssertFalse(t, got)
}

func TestCompiledShouldReport500(t *testing.T) {
	res := http.Response{Code: 500}
	checker := Compile("code = 500")

	got := checker(res)

	testutils.AssertTrue(t, got)
}

func TestCompiledShouldReport403(t *testing.T) {
	res := http.Response{Code: 403}
	checker := Compile("code != 200")

	got := checker(res)

	testutils.AssertTrue(t, got)
}

func TestCompiledShouldMatchLength(t *testing.T) {
	res := http.Response{Code: 200, Length: 1500}
	checker := Compile("size = 1500")

	got := checker(res)

	testutils.AssertTrue(t, got)
}

func TestCompiledShouldMatchText(t *testing.T) {
	res := http.Response{Code: 200, Raw: []byte("hello foo bar")}
	checker := Compile("text = foo")

	got := checker(res)

	testutils.AssertTrue(t, got)
}

func TestCompiledShouldMatchBothConditions(t *testing.T) {
	res := http.Response{Code: 200, Length: 1000}
	checker := Compile("code = 200 and size = 1500")

	got := checker(res)

	testutils.AssertFalse(t, got)
}

func TestCompiledShouldMatchOneCondition(t *testing.T) {
	res := http.Response{Code: 200, Length: 1000}
	checker := Compile("code = 200 or size = 1500")

	got := checker(res)

	testutils.AssertTrue(t, got)
}
