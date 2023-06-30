package matchlang

import (
	"github.com/kamil-s-solecki/haze/testutils"
	"testing"
)

func TestASingleComparisonIsValid(t *testing.T) {
	ok, err := Validate("code = 200")

	testutils.AssertTrue(t, ok)
	testutils.AssertNil(t, err)
}

func TestAnEmptyStringIsInvalid(t *testing.T) {
	ok, err := Validate("")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "The expression cannot be empty!")
}

func TestABlankStringIsInvalid(t *testing.T) {
	ok, err := Validate("      ")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "The expression cannot be empty!")
}

func TestExpressionThatDoesNotStartWithAndIdentifierIsInvalid(t *testing.T) {
	ok, err := Validate("foo = 200")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "foo is not a valid identifier!")
}
