package calculator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculator(t *testing.T) {
	const MaxRecipientCnt = 1

	restrictions := map[string][]string{
		"Дима":   {"Ира"},
		"Ира":    {"Дима"},
		"Никита": {"Люба"},
		"Люба":   {"Никита"},
		"Миша":   {"Настя"},
		"Настя":  {"Миша"},
	}

	Convey("test calculator", t, func() {
		c := NewCalculator(restrictions)
		res := c.CalculateRecipient()

		fmt.Printf("%+v\n", res)

		recipientCounts := make(map[string]int, len(restrictions))
		for p := range restrictions {
			recipientCounts[p] = 0
		}

		for sender, recipient := range res {
			recipientCounts[string(recipient)]++
			So(sender, ShouldNotEqual, recipient)
			So(recipient, ShouldNotBeIn, restrictions[string(sender)])
		}

		for p := range restrictions {
			So(recipientCounts[p], ShouldEqual, MaxRecipientCnt)
		}
	})
}
