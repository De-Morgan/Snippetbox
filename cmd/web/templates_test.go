package main

import (
	"testing"
	"time"
)

func Test_humanDate(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "utc",
			args: args{
				t: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			},
			want: "17 Dec 2020 at 10:00",
		},
		{name: "Empty",
			args: args{
				t: time.Time{},
			},
			want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := humanDate(tt.args.t); got != tt.want {
				t.Errorf("humanDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
