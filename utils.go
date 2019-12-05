package main

import (
	"sort"
	"strings"
)

type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func GetNamespaceParts(ns string) []string {
	parts := strings.Split(ns, "/")
	subparts := []string{ns}
	for index := 0; index < len(parts); index++ {
		part := strings.Join(parts[0:index], "/")
		if len(part) == 0 {
			continue
		}
		subparts = append(subparts, part)
	}
	sort.Sort(byLength(subparts))
	return subparts
}

func EqualInt(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
