package main

import (
	"io/ioutil"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type Config struct {
	GroupID int                 `yaml:"groupID"`
	Options map[string]*Options `yaml:"options"`
}

type Options struct {
	AllowedApprovers      gitlab.ChangeAllowedApproversOptions      `yaml:"allowedApprovers"`
	ApprovalConfiguration gitlab.ChangeApprovalConfigurationOptions `yaml:"approvalConfiguration"`
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
