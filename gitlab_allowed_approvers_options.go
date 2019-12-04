package main

import (
	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasAllowedApproversChanges(project *gitlab.Project, opts *Options) (changed bool) {
	aa := opts.AllowedApprovers
	approvals, _, err := d.GitLab.Projects.GetApprovalConfiguration(project.ID)
	if err != nil {
		return false
	}
	if aa.ApproverIDs != nil {
		if len(aa.ApproverIDs) != len(approvals.Approvers) {
			logrus.Infof(
				"Approvers: %d != %d",
				len(aa.ApproverIDs),
				len(approvals.Approvers),
			)
			changed = true
		}
	}
	return changed
}

func (d *Dredd) FixAllowedApprovers(project *gitlab.Project, opts *Options) error {
	_, _, err := d.GitLab.Projects.ChangeAllowedApprovers(project.ID, &opts.AllowedApprovers)
	return err
}
