package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasProjectApprovalRuleChanges(project *gitlab.Project, opts *Options) bool {
	if opts.ApprovalRule == nil {
		return false
	}

	rules, _, err := d.GitLab.Projects.GetProjectApprovalRules(project.ID)
	if err != nil {
		return false
	}

	for _, rule := range rules {
		if rule.Name != *opts.ApprovalRule.Name {
			continue
		}
		if rule.ApprovalsRequired != *opts.ApprovalRule.ApprovalsRequired {
			logrus.Infof("Project approval rule %s approvals required: %d != %d", rule.Name, rule.ApprovalsRequired, *opts.ApprovalRule.ApprovalsRequired)
			return true
		}
		var userIDs []int
		for _, user := range rule.Users {
			userIDs = append(userIDs, user.ID)
		}
		if !EqualInt(userIDs, opts.ApprovalRule.UserIDs) {
			logrus.Infof("Project approval rule %s user IDs: %v != %v", rule.Name, userIDs, opts.ApprovalRule.UserIDs)
			return true
		}
		return false
	}

	logrus.Infof("Project approval rule %s not found", *opts.ApprovalRule.Name)
	return true
}

func (d *Dredd) FixProjectApprovalRule(project *gitlab.Project, opts *Options) error {
	rules, _, err := d.GitLab.Projects.GetProjectApprovalRules(project.ID)
	if err != nil {
		return nil
	}
	for _, rule := range rules {
		if rule.Name == *opts.ApprovalRule.Name {
			updateOpts := &gitlab.UpdateProjectLevelRuleOptions{
				ApprovalsRequired: opts.ApprovalRule.ApprovalsRequired,
				UserIDs:           opts.ApprovalRule.UserIDs,
			}
			_, _, err := d.GitLab.Projects.UpdateProjectApprovalRule(project.ID, rule.ID, updateOpts)
			return err
		}
	}
	_, _, err = d.GitLab.Projects.CreateProjectApprovalRule(project.ID, opts.ApprovalRule)
	return err
}
