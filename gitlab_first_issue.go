package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

const (
	firstIssueLabel = "good-first-issue"
)

func (d *Dredd) HasFirstIssueChanges(project *gitlab.Project, opts *Options) bool {
	if opts.FirstIssue == nil {
		return false
	}

	listOpts := &gitlab.ListProjectIssuesOptions{
		Labels: gitlab.Labels{firstIssueLabel},
	}

	issues, _, err := d.GitLab.Issues.ListProjectIssues(project.ID, listOpts)
	if err != nil {
		return false
	}

	if len(issues) > 0 {
		return false
	}

	logrus.Infof("Create %s issues: %d != 1", firstIssueLabel, len(issues))
	return true
}

func (d *Dredd) CreateFirstIssue(project *gitlab.Project, opts *Options) error {
	createOpts := &gitlab.CreateIssueOptions{
		Title:       opts.FirstIssue.Title,
		Description: opts.FirstIssue.Description,
		Labels:      &gitlab.Labels{firstIssueLabel},
		AssigneeIDs: []int{project.CreatorID},
	}
	_, _, err := d.GitLab.Issues.CreateIssue(project.ID, createOpts)
	return err
}
