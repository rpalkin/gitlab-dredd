package main

import (
	"io/ioutil"
	"os"

	gitlab "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

const (
	defaultBaseURL = "https://gitlab.com/"
)

type Config struct {
	GitLabEndpoint string              `yaml:"gitlabEndpoint"`
	GitLabToken    string              `yaml:"gitlabToken"`
	GroupIDs       []int               `yaml:"groupIDs"`
	ListenAddress  string              `yaml:"listenAddress"`
	Options        map[string]*Options `yaml:"options"`
}

type Options struct {
	AllowedApprovers      *gitlab.ChangeAllowedApproversOptions      `yaml:"allowedApprovers"`
	ApprovalConfiguration *gitlab.ChangeApprovalConfigurationOptions `yaml:"approvalConfiguration"`
	ApprovalRule          *gitlab.CreateProjectLevelRuleOptions      `yaml:"approvalRule"`
	FirstIssue            *gitlab.CreateIssueOptions                 `yaml:"firstIssue"`
	JiraIntegration       *gitlab.SetJiraServiceOptions              `yaml:"jiraIntegration"`
	ProjectOptions        *gitlab.EditProjectOptions                 `yaml:"projectOptions"`
	PushRules             *gitlab.EditProjectPushRuleOptions         `yaml:"pushRules"`
	RepositoryBranches    []*gitlab.ProtectRepositoryBranchesOptions `yaml:"repositoryBranches"`
}

func LoadFromFile(filename string) (*Config, error) {
	config := &Config{}
	confContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	confContent = []byte(os.ExpandEnv(string(confContent)))
	err = yaml.Unmarshal(confContent, &config)
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
	if len(config.ListenAddress) == 0 {
		config.ListenAddress = ":8080"
	}
	return config, nil
}

func (c *Config) OptionsByPath(ns string) *Options {
	parts := GetNamespaceParts(ns)
	for _, part := range parts {
		if opt, ok := c.Options[part]; ok {
			return opt
		}
	}
	return nil
}
