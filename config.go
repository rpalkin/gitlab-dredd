package main

import (
	"io/ioutil"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GroupIDs []int               `yaml:"groupIDs"`
	Options  map[string]*Options `yaml:"options"`
}

type Options struct {
	AllowedApprovers      gitlab.ChangeAllowedApproversOptions      `yaml:"allowedApprovers"`
	ApprovalConfiguration gitlab.ChangeApprovalConfigurationOptions `yaml:"approvalConfiguration"`
	ProjectOptions        gitlab.EditProjectOptions                 `yaml:"projectOptions"`
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
