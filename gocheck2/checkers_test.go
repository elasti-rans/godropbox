package gocheck2

import (
	stdlibErrors "errors"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"

	"github.com/dropbox/godropbox/errors"
)

// Hook up gocheck into go test runner
func Test(t *testing.T) {
	TestingT(t)
}

type CheckersSuite struct{}

var _ = Suite(&CheckersSuite{})

func (s *CheckersSuite) SetUpTest(c *C) {
}

func test(c *C, ch Checker,
	expectedResult bool, expectedErr string, params ...interface{}) {

	actualResult, actualErr := ch.Check(params, nil)
	if actualResult != expectedResult || actualErr != expectedErr {
		c.Fatalf(
			"Check returned (%#v, %#v) rather than (%#v, %#v)",
			actualResult, actualErr, expectedResult, expectedErr)
	}
}

func (s *CheckersSuite) TestHasKey(c *C) {
	test(c, HasKey, true, "", map[string]int{"foo": 1}, "foo")
	test(c, HasKey, false, "", map[string]int{"foo": 1}, "bar")
	test(c, HasKey, true, "", map[int][]byte{10: nil}, 10)

	test(c, HasKey, false, "First argument to HasKey must be a map", nil, "bar")
	test(c, HasKey,
		false, "Second argument must be assignable to the map key type",
		map[string]int{"foo": 1}, 10)
}

func (s *CheckersSuite) TestNoErr(c *C) {
	// Test the true/false behavior.
	test(c, NoErr, true, "", nil)
	test(c, NoErr, true, "", 3)
	test(c, NoErr, true, "", error(nil))
	test(c, NoErr, false, "", stdlibErrors.New("message"))
	test(c, NoErr, false, "", errors.New("message"))

	// Test the message behavior.
	params := []interface{}{errors.New("1\n2\n3")}
	text := params[0].(error).Error()
	NoErr.Check(params, nil)
	c.Assert(fmt.Sprintf("%#v", params[0]), Equals, "\n" + text)
}

func (s *CheckersSuite) TestErrorMatches(c *C) {
	test(c, MultilineErrorMatches, true, "",
		errors.Newf("Oh damn, this stinks"), "stinks")
	test(c, MultilineErrorMatches, false, "",
		errors.Newf("Oh damn, this stinks"), "skinks")
}
