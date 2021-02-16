package gene

import (
	"reflect"
	"testing"

	"github.com/misterikkit/automata/life/game"
)

func TestGene_String(t *testing.T) {
	type fields struct {
		Alive [9]game.Cell
		Dead  [9]game.Cell
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "all zero",
			// fields: {},
			want: "000000000000000000",
		},
		{
			name: "all one",
			fields: fields{
				Alive: [9]game.Cell{true, true, true, true, true, true, true, true, true},
				Dead:  [9]game.Cell{true, true, true, true, true, true, true, true, true},
			},
			want: "111111111111111111",
		},
		{
			name: "alternate",
			fields: fields{
				Alive: [9]game.Cell{true, false, true, false, true, false, true, false, true},
				Dead:  [9]game.Cell{false, true, false, true, false, true, false, true, false},
			},
			want: "101010101010101010",
		},
		{
			name: "AB",
			fields: fields{
				Alive: [9]game.Cell{true, true, true, true, true, true, true, true, true},
				Dead:  [9]game.Cell{false, false, false, false, false, false, false, false, false},
			},
			want: "111111111000000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Gene{
				Alive: tt.fields.Alive,
				Dead:  tt.fields.Dead,
			}
			if got := g.String(); got != tt.want {
				t.Errorf("Gene.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromString(t *testing.T) {

	tests := []struct {
		name string
		arg  string
		want Gene
	}{
		{
			name: "all zero",
			arg:  "000000000000000000",
			want: Gene{
				Alive: [9]game.Cell{false, false, false, false, false, false, false, false, false},
				Dead:  [9]game.Cell{false, false, false, false, false, false, false, false, false},
			},
		},
		{
			name: "all one",
			arg:  "111111111111111111",
			want: Gene{
				Alive: [9]game.Cell{true, true, true, true, true, true, true, true, true},
				Dead:  [9]game.Cell{true, true, true, true, true, true, true, true, true},
			},
		},
		{
			name: "alternate",
			arg:  "101010101010101010",
			want: Gene{
				Alive: [9]game.Cell{true, false, true, false, true, false, true, false, true},
				Dead:  [9]game.Cell{false, true, false, true, false, true, false, true, false},
			},
		},
		{
			name: "AB",
			arg:  "111111111000000000",
			want: Gene{
				Alive: [9]game.Cell{true, true, true, true, true, true, true, true, true},
				Dead:  [9]game.Cell{false, false, false, false, false, false, false, false, false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromString(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
