package cycle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet_Add(t *testing.T) {
	type fields struct {
		head *node
	}
	type args struct {
		v Val
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Set{
				head: tt.fields.head,
			}
			s.Add(tt.args.v)
		})
	}
}

func TestSet_toSlice(t *testing.T) {
	require := require.New(t)
	s := Set{}
	sl := s.toSlice()
	require.Empty(sl, "Empty set should return empty slice")
	s.Add(12)
	sl = s.toSlice()
	require.ElementsMatch(sl, []Val{12}, "Set with one element should return that element")
	for i := 0; i < 5; i++ {
		s.Add(i)
	}
	sl = s.toSlice()
	require.ElementsMatch(sl, []Val{0, 1, 2, 3, 4, 12})
	// duplicates allowed
	s.Add(0)
	sl = s.toSlice()
	require.ElementsMatch(sl, []Val{0, 0, 1, 2, 3, 4, 12})
}

func TestSet_Merge(t *testing.T) {
	require := require.New(t)
	tests := []struct {
		name   string
		s1, s2 *Set
		want   []Val
	}{
		{
			name: "both empty",
			s1:   new(Set),
			s2:   new(Set),
			want: nil,
		},
		{
			name: "left empty",
			s1:   new(Set),
			s2:   fromVals(1),
			want: []Val{1},
		},
		{
			name: "right empty",
			s1:   fromVals(1),
			s2:   new(Set),
			want: []Val{1},
		},
		{
			name: "left empty/multi",
			s1:   new(Set),
			s2:   fromVals(1, 2, 3),
			want: []Val{1, 2, 3},
		},
		{
			name: "right empty/multi",
			s1:   fromVals(1, 2, 3),
			s2:   new(Set),
			want: []Val{1, 2, 3},
		},

		{
			name: "single-single",
			s1:   fromVals(1),
			s2:   fromVals(2),
			want: []Val{1, 2},
		},
		{
			name: "single-multi",
			s1:   fromVals(1),
			s2:   fromVals(2, 3, 4),
			want: []Val{1, 2, 3, 4},
		},
		{
			name: "multi-single",
			s1:   fromVals(1, 2, 3),
			s2:   fromVals(4),
			want: []Val{1, 2, 3, 4},
		},
		{
			name: "multi-multi",
			s1:   fromVals(1, 2, 3),
			s2:   fromVals(4, 5, 6),
			want: []Val{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s1.Merge(tt.s2)
			require.ElementsMatch(tt.s1.toSlice(), tt.want, "Left set does not match expected")
			require.ElementsMatch(tt.s2.toSlice(), tt.want, "Right set does not match expected")
		})
	}
}

func fromVals(vs ...Val) *Set {
	s := &Set{}
	for _, v := range vs {
		s.Add(v)
	}
	return s
}
