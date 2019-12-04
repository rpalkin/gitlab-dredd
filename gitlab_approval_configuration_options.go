package main

import (
	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasApprovalConfigurationChanges(project *gitlab.Project, opts *Options) (changed bool) {
	ac := opts.ApprovalConfiguration
	approvals, _, err := d.GitLab.Projects.GetApprovalConfiguration(project.ID)
	if err != nil {
		return false
	}
	if ac.ApprovalsBeforeMerge != nil {
		if *ac.ApprovalsBeforeMerge != approvals.ApprovalsBeforeMerge {
			logrus.Infof(
				"Approvals Before Merge: %d != %d",
				*ac.ApprovalsBeforeMerge,
				approvals.ApprovalsBeforeMerge,
			)
			changed = true
		}
	}
	if ac.ResetApprovalsOnPush != nil {
		if *ac.ResetApprovalsOnPush != approvals.ResetApprovalsOnPush {
			logrus.Infof(
				"Reset Approvals On Push: %v != %v",
				*ac.ResetApprovalsOnPush,
				approvals.ResetApprovalsOnPush,
			)
			changed = true
		}
	}
	if ac.DisableOverridingApproversPerMergeRequest != nil {
		if *ac.DisableOverridingApproversPerMergeRequest != approvals.DisableOverridingApproversPerMergeRequest {
			logrus.Infof(
				"Disable Overriding Approvers Per Merge Request: %v != %v",
				*ac.DisableOverridingApproversPerMergeRequest,
				approvals.DisableOverridingApproversPerMergeRequest,
			)
			changed = true
		}
	}
	if ac.MergeRequestsAuthorApproval != nil {
		if *ac.MergeRequestsAuthorApproval != approvals.MergeRequestsAuthorApproval {
			logrus.Infof(
				"Merge Requests Author Approval: %v != %v",
				*ac.MergeRequestsAuthorApproval,
				approvals.MergeRequestsAuthorApproval,
			)
			changed = true
		}
	}
	return changed
}

func (d *Dredd) FixApprovalConfiguration(project *gitlab.Project, opts *Options) error {
	_, _, err := d.GitLab.Projects.ChangeApprovalConfiguration(project.ID, &opts.ApprovalConfiguration)
	return err
}
