package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasProjectPushRulesChanges(project *gitlab.Project, opts *Options) (changed bool) {
	if opts.PushRules == nil {
		return false
	}

	po := opts.PushRules
	pushRules, _, err := d.GitLab.Projects.GetProjectPushRules(project.ID)
	if err != nil {
		return false
	}

	if po.BranchNameRegex != nil {
		if *po.BranchNameRegex != pushRules.BranchNameRegex {
			logrus.Infof("Branch name regexp: %s != %s",
				*po.BranchNameRegex,
				pushRules.BranchNameRegex,
			)
			changed = true
		}
	}

	if po.CommitMessageRegex != nil {
		if *po.CommitMessageRegex != pushRules.CommitMessageRegex {
			logrus.Infof("Commit message regexp: %s != %s",
				*po.CommitMessageRegex,
				pushRules.CommitMessageRegex,
			)
			changed = true
		}
	}

	return changed
}

func (d *Dredd) FixProjectPushRules(project *gitlab.Project, opts *Options) error {
	_, _, err := d.GitLab.Projects.EditProjectPushRule(project.ID, opts.PushRules)
	return err
}
