package region

import (
	"math/rand"
	"testing"
	"time"

	"github.com/misterikkit/automata/maze/game"
	"github.com/misterikkit/automata/maze/gene"
	"github.com/misterikkit/automata/maze/tui"
)

func TestMap(t *testing.T) {
	rand.Seed(time.Now().Unix())
	gn := gene.Random()
	t.Logf("Gene: %v", gn)
	rule := gn.AsRule()
	g := game.New(100, 100)
	g = g.Next(game.Random)
	for i := 0; i < 80; i++ {
		g = g.Next(rule)
	}
	t.Logf("Game state:\n%v", tui.Fmt(g))
	mapped := Map(g)
	t.Logf("Game map:\n%v", mapped)
}

func Test_follow(t *testing.T) {
	type args struct {
		m map[int]int
		k int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty",
			args: args{nil, 20},
			want: 20,
		},
		{
			name: "miss",
			args: args{map[int]int{1: 2}, 5},
			want: 5,
		},
		{
			name: "simple",
			args: args{map[int]int{11: 5, 5: 3}, 11},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := follow(tt.args.m, tt.args.k); got != tt.want {
				t.Errorf("follow() = %v, want %v", got, tt.want)
			}
		})
	}
}
