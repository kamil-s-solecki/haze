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
