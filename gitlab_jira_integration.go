package main

import (
	"github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasJiraIntegrationChanges(project *gitlab.Project, opts *Options) bool {
	if opts.JiraIntegration == nil {
		return false
	}

	desiredOpts := opts.JiraIntegration

	jira, _, err := d.GitLab.Services.GetJiraService(project.ID)
	if err != nil {
		return false
	}
	return *desiredOpts.Active != jira.Active ||
		*desiredOpts.MergeRequestsEvents != jira.MergeRequestsEvents ||
		*desiredOpts.CommentOnEventEnabled != jira.CommentOnEventEnabled ||
		*desiredOpts.CommitEvents != jira.CommitEvents ||
		*desiredOpts.URL != jira.Properties.URL ||
		*desiredOpts.Username != jira.Properties.Username ||
		*desiredOpts.Password != jira.Properties.Password
}

func (d *Dredd) SetJiraIntegration(project *gitlab.Project, opts *Options) error {
	_, err := d.GitLab.Services.SetJiraService(project.ID, opts.JiraIntegration)
	return err
}
