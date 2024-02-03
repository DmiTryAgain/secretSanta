package calculator

import "testing"

func TestParticipantIDs_isEqualValues(t *testing.T) {
	type args struct {
		ps ParticipantIDs
	}
	tests := []struct {
		name string
		pp   ParticipantIDs
		args args
		want bool
	}{
		{
			name: "not equal with different len",
			pp: ParticipantIDs{
				"1",
				"2",
				"3",
			},
			args: args{
				ps: ParticipantIDs{
					"1",
					"2",
				},
			},
			want: false,
		},
		{
			name: "not equal with the same len",
			pp: ParticipantIDs{
				"1",
				"2",
				"3",
			},
			args: args{
				ps: ParticipantIDs{
					"1",
					"2",
					"4",
				},
			},
			want: false,
		},
		{
			name: "equal with the same len and order",
			pp: ParticipantIDs{
				"1",
				"2",
				"3",
			},
			args: args{
				ps: ParticipantIDs{
					"1",
					"2",
					"3",
				},
			},
			want: true,
		},
		{
			name: "equal with the same len and different order 1",
			pp: ParticipantIDs{
				"1",
				"2",
				"3",
			},
			args: args{
				ps: ParticipantIDs{
					"3",
					"1",
					"2",
				},
			},
			want: true,
		},
		{
			name: "equal with the same len and different order 2",
			pp: ParticipantIDs{
				"2",
				"1",
				"3",
			},
			args: args{
				ps: ParticipantIDs{
					"1",
					"2",
					"3",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pp.isEqualValues(tt.args.ps); got != tt.want {
				t.Errorf("isEqualValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
