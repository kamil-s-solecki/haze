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

func TestEpressionWithoutAnOperatorIsInvalid(t *testing.T) {
	ok, err := Validate("code foo 200")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "foo is not a valid operator!")
}

func TestEpressionWithoutALiteralIsInvalid(t *testing.T) {
	ok, err := Validate("code = code")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "code is not a valid literal!")
}

func TestEpressionWithInvalidOperatorIsInvalid(t *testing.T) {
	ok, err := Validate("code = 200 foo text = bar")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "foo is not a valid logical operator!")
}

func TestExpressionWithInvalidSecondComparisonIsInvalid(t *testing.T) {
	ok, err := Validate("code = 200 and foo = bar")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "foo is not a valid identifier!")
}

func TestExpressionWithSecondComparisonWithoutALiteralIsInvalid(t *testing.T) {
	ok, err := Validate("code = 200 and code =")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "Expected a literal after '='!")
}

func TestExpressionWithSecondComparisonWithoutAnOperatorIsInvalid(t *testing.T) {
	ok, err := Validate("code = 200 and code")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "Expected an operator after 'code'!")
}

func TestExpressionWithoutSecondComparisonIsInvalid(t *testing.T) {
	ok, err := Validate("code = 200 and")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "Expected a comparison after 'and'!")
}

func TestExpressionWithoutThirdComparisonIsInvalid(t *testing.T) {
	ok, err := Validate("code = 200 and code = 200 and")

	testutils.AssertFalse(t, ok)
	testutils.AssertErrorEquals(t, err, "Expected a comparison after 'and'!")
}
