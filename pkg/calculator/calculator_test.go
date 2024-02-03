package calculator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculator(t *testing.T) {
	const MaxRecipientCnt = 1
	type args struct {
		restrictions map[string][]string
	}

	var tests = []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "2 participants",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {},
					"Света": {},
				},
			},
			wantErr: ErrNotEnoughParticipants,
		},
		{
			name: "3 participants with invalid restrictions",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света"},
					"Света": {"Петя"},
					"Паша":  {},
				},
			},
			wantErr: ErrIncorrectRestrictions,
		},
		{
			name: "3 participants with valid restrictions",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света"},
					"Света": {"Паша"},
					"Паша":  {"Петя"},
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants with valid restrictions 1",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света"},
					"Света": {"Паша"},
					"Паша":  {"Галя"},
					"Галя":  {"Петя"},
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants with valid restrictions 2",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света"},
					"Света": {"Петя"},
					"Паша":  {"Галя"},
					"Галя":  {"Паша"},
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants with valid restrictions 3",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света", "Паша"},
					"Света": {"Петя"},
					"Паша":  {"Галя"},
					"Галя":  {"Паша"},
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants with valid restrictions 4",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света", "Паша", "Петя", "Оксана"},
					"Света": {"Петя"},
					"Паша":  {"Галя"},
					"Галя":  {"Паша"},
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants with valid restrictions 5",
			args: args{
				restrictions: map[string][]string{
					"Петя":  {"Света", "Паша"}, //Гале
					"Света": {"Петя", "Галя"},  //Паше
					"Паша":  {"Галя", "Петя"},  //Свете
					"Галя":  {"Паша", "Света"}, //Пете
				},
			},
			wantErr: nil,
		},
		{
			name: "with one restrict for each participant",
			args: args{
				restrictions: map[string][]string{
					"Дима":   {"Ира"},
					"Ира":    {"Дима"},
					"Никита": {"Люба"},
					"Люба":   {"Никита"},
					"Миша":   {"Настя"},
					"Настя":  {"Миша"},
				},
			},
			wantErr: nil,
		},
	}

	Convey("test calculator", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				arg := tt.args.restrictions
				c, err := NewCalculator(arg)
				So(tt.wantErr, ShouldBeError, err)
				if err != nil {
					return
				}

				res := c.CalculateRecipient()
				fmt.Printf("%+v\n", res)

				recipientCounts := make(map[string]int, len(arg))
				for p := range arg {
					recipientCounts[p] = 0
				}

				for sender, recipient := range res {
					recipientCounts[string(recipient)]++
					So(sender, ShouldNotEqual, recipient)
					So(recipient, ShouldNotBeIn, arg[string(sender)])
				}

				for p := range arg {
					So(recipientCounts[p], ShouldEqual, MaxRecipientCnt)
				}
			})
		}
	})
}
