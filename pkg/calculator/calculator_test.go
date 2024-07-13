package calculator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalculator(t *testing.T) {
	const MaxRecipientCnt = 1
	type args struct {
		participants map[string]string
		restrictions map[string][]string
	}

	var tests = []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "with nil participants",
			args: args{
				participants: nil,
				restrictions: nil,
			},
			wantErr: ErrNotEnoughParticipants,
		},
		{
			name: "without participants",
			args: args{
				participants: map[string]string{},
				restrictions: nil,
			},
			wantErr: ErrNotEnoughParticipants,
		},
		{
			name: "2 participants",
			args: args{
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
				},
			},
			wantErr: ErrNotEnoughParticipants,
		},
		{
			name: "3 participants with invalid restrictions",
			args: args{
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
				},
				restrictions: map[string][]string{
					"Петя":  {"Света"},
					"Света": {"Петя"},
				},
			},
			wantErr: ErrIncorrectRestrictions,
		},
		{
			name: "3 participants with no available for one of them",
			args: args{
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
				},
				restrictions: map[string][]string{
					"Петя":  {"Света", "Паша"},
					"Света": {"Петя"},
				},
			},
			wantErr: ErrIncorrectRestrictions,
		},
		{
			name: "3 participants with valid restrictions",
			args: args{
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
				},
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
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
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
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
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
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
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
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
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
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
				restrictions: map[string][]string{
					"Петя":  {"Света", "Паша"}, // Гале
					"Света": {"Петя", "Галя"},  // Паше
					"Паша":  {"Галя", "Петя"},  // Свете
					"Галя":  {"Паша", "Света"}, // Пете
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants with non-existed restrictions 5",
			args: args{
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
				restrictions: map[string][]string{
					"Петя":  {"Света", "Паша"},
					"Света": {"Петя", "Галя"},
					"Толик": {"Петя", "Света", "Паша", "Галя"},
				},
			},
			wantErr: nil,
		},
		{
			name: "4 participants without restrictions",
			args: args{
				participants: map[string]string{
					"Петя":  "Петя",
					"Света": "Света",
					"Паша":  "Паша",
					"Галя":  "Галя",
				},
				restrictions: map[string][]string{},
			},
			wantErr: nil,
		},
		{
			name: "with one restrict for each participant",
			args: args{
				participants: map[string]string{
					"Дима":   "Дима",
					"Ира":    "Ира",
					"Никита": "Никита",
					"Люба":   "Люба",
					"Миша":   "Миша",
					"Настя":  "Настя",
				},
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
				participants := tt.args.participants
				restrictions := tt.args.restrictions
				c, err := NewCalculator(participants, restrictions)
				if tt.wantErr != nil {
					So(err, ShouldNotBeNil)
					So(err, ShouldWrap, tt.wantErr)
					return
				} else if err != nil {
					t.Fatal(err)
				}

				res := c.CalculateRecipient()
				fmt.Printf("%+v\n", res)

				recipientCounts := make(map[string]int, len(participants))
				for p := range participants {
					recipientCounts[p] = 0
				}

				for sender, recipient := range res {
					recipientCounts[string(recipient)]++
					So(sender, ShouldNotEqual, recipient)
					So(recipient, ShouldNotBeIn, restrictions[string(sender)])
				}

				for p := range participants {
					So(recipientCounts[p], ShouldEqual, MaxRecipientCnt)
				}
			})
		}
	})
}
