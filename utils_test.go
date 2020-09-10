package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamespaceParts(t *testing.T) {
	tests := []struct {
		path string
		res  []string
	}{
		{
			path: "group/project",
			res:  []string{"group/project", "group"},
		},
		{
			path: "group/subgroup/project",
			res:  []string{"group/subgroup/project", "group/subgroup", "group"},
		},
		{
			path: "group/subgroup/subgroup/project",
			res:  []string{"group/subgroup/subgroup/project", "group/subgroup/subgroup", "group/subgroup", "group"},
		},
	}
	for _, test := range tests {
		res := GetNamespaceParts(test.path)
		if !reflect.DeepEqual(res, test.res) {
			t.Errorf("%v != %v", res, test.res)
		}
	}
}

func TestEqualInt(t *testing.T) {
	tests := []struct {
		in1 []int
		in2 []int
		eq  bool
	}{
		{
			in1: []int{1, 2, 3},
			in2: []int{},
			eq:  false,
		},
		{
			in1: []int{1, 2, 3},
			in2: []int{1, 2, 3},
			eq:  true,
		},
		{
			in1: []int{1, 2, 3},
			in2: []int{1, 3, 4},
			eq:  false,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.eq, EqualInt(test.in1, test.in2))
	}
}
