package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestEller(t *testing.T) {
	rand.Seed(time.Now().Unix())

	maze := buildEllerMaze(10, 10)
	t.Log(maze)
}

func Test_setID_String(t *testing.T) {
	tests := []struct {
		name string
		i    setID
		want string
	}{
		{i: 0, want: "0", name: "zero"},
		{i: 1, want: "1", name: "one"},
		{i: 10, want: "A", name: "alpha"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("setID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
