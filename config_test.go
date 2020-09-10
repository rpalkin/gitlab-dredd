package main

import (
	"os"
	"testing"

	"github.com/xanzy/go-gitlab"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromFile(t *testing.T) {
	_, err := LoadFromFile("test_data/not_found.yaml")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "no such file or directory")
	}

	_, err = LoadFromFile("test_data/invalid-syntax.yaml")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "yaml")
	}

	os.Setenv("GITLAB_ENDPOINT", "endpoint")
	os.Setenv("GITLAB_TOKEN", "token")

	config, err := LoadFromFile("test_data/empty-config.yaml")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, &Config{
		ListenAddress:  ":8080",
		GitLabEndpoint: "endpoint",
		GitLabToken:    "token",
		GroupIDs:       []int{89},
	}, config)

	os.Unsetenv("GITLAB_ENDPOINT")

	config, err = LoadFromFile("test_data/empty-config.yaml")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, &Config{
		ListenAddress:  ":8080",
		GitLabEndpoint: "https://gitlab.com/",
		GitLabToken:    "token",
		GroupIDs:       []int{89},
	}, config)

	config, err = LoadFromFile("test_data/nonempty-config.yaml")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, &Config{
		ListenAddress:  "0.0.0.0:8080",
		GitLabEndpoint: "http://gitlab.local",
		GitLabToken:    "ABCD",
		GroupIDs:       []int{89},
	}, config)
}

func TestConfig_OptionsByPath(t *testing.T) {
	c := &Config{
		Options: map[string]*Options{
			"abcd": {
				AllowedApprovers: &gitlab.ChangeAllowedApproversOptions{
					ApproverIDs: []int{1},
				},
			},
		},
	}
	assert.Nil(t, c.OptionsByPath(""))
	assert.Nil(t, c.OptionsByPath("foobar"))
	assert.Equal(t, &Options{
		AllowedApprovers: &gitlab.ChangeAllowedApproversOptions{
			ApproverIDs: []int{1},
		},
	}, c.OptionsByPath("abcd"))
}
