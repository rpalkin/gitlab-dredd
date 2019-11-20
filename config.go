package main

import (
	"io/ioutil"
	"os"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

const (
	defaultBaseURL = "https://gitlab.com/"
)

type Config struct {
	ListenAddress  string              `yaml:"listenAddress"`
	GitLabEndpoint string              `yaml:"gitlabEndpoint"`
	GitLabToken    string              `yaml:"gitlabToken"`
	GroupIDs       []int               `yaml:"groupIDs"`
	Options        map[string]*Options `yaml:"options"`
}

type Options struct {
	AllowedApprovers      gitlab.ChangeAllowedApproversOptions       `yaml:"allowedApprovers"`
	ApprovalConfiguration gitlab.ChangeApprovalConfigurationOptions  `yaml:"approvalConfiguration"`
	ProjectOptions        gitlab.EditProjectOptions                  `yaml:"projectOptions"`
	RepositoryBranches    []*gitlab.ProtectRepositoryBranchesOptions `yaml:"repositoryBranches"`
}

func LoadFromFile(filename string) (*Config, error) {
	config := &Config{}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	if len(config.GitLabEndpoint) == 0 {
		config.GitLabEndpoint = os.Getenv("GITLAB_ENDPOINT")
		if len(config.GitLabEndpoint) == 0 {
			config.GitLabEndpoint = defaultBaseURL
		}
	}
	if len(config.GitLabToken) == 0 {
		config.GitLabToken = os.Getenv("GITLAB_TOKEN")
	}
	return config, nil
}

func (c *Config) OptionsByPath(ns string) *Options {
	parts := GetNamespaceParts(ns)
	for _, part := range parts {
		opt, ok := c.Options[part]
		if !ok {
			continue
		}
		return opt
	}
	return nil
}
