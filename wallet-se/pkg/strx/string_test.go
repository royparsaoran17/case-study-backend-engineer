package strx

import (
	"fmt"
	"testing"
)

func TestFirstToUpper(t *testing.T) {
	cases := []struct {
		Input    string
		Expected string
	}{
		{
			Input:    "",
			Expected: "",
		},

		{
			Input:    "testing",
			Expected: "Testing",
		},

		{
			Input:    "testing lagi",
			Expected: "Testing lagi",
		},

		{
			Input:    "TESTING",
			Expected: "TESTING",
		},
	}

	for _, v := range cases {

		actual := FirstToUpper(v.Input)
		if actual == v.Expected {
			t.Logf(fmt.Sprintf("string first to upper, expected %s, actual %s", v.Expected, actual))
		} else {

			t.Errorf(fmt.Sprintf("string first to upper, expected %s, actual %s", v.Expected, actual))
		}
	}

}
