package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/xanzy/go-gitlab"
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

func GetStdinHookPayload() (*gitlab.HookEvent, error) {
	reader := bufio.NewReader(os.Stdin)
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	hook := &gitlab.HookEvent{}
	err = json.Unmarshal(b, &hook)
	if err != nil {
		return nil, err
	}
	return hook, nil
}
