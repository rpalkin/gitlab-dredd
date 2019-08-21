package main

import (
	"reflect"
	"testing"
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
