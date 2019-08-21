package main

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type Dredd struct {
	GitLab *gitlab.Client
	Config *Config
	DryRun bool
}

func (d *Dredd) Run() error {
	logrus.Info("Requesting projects list...")
	projects, err := d.GetProjects(d.Config.GroupID)
	if err != nil {
		return err
	}
	err = d.processProjects(projects)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dredd) processProjects(projects []*gitlab.Project) error {
	logrus.Infof("Found projects: %d", len(projects))
	for _, project := range projects {
		err := d.processProject(project)
		if err != nil {
			logrus.Error(err)
			continue
		}
	}
	return nil
}

func (d *Dredd) processProject(project *gitlab.Project) error {
	logrus.Infof("Processing %s project...", project.Name)
	approvals, _, err := d.GitLab.Projects.GetApprovalConfiguration(project.ID)
	if err != nil {
		return err
	}
	opts := d.Config.OptionsByPath(project.PathWithNamespace)
	if opts == nil {
		return nil
	}
	if !d.hasProjectApprovalsChanges(opts, approvals) {
		return nil
	}
	if d.DryRun {
		return nil
	}
	_, _, err = d.GitLab.Projects.ChangeAllowedApprovers(project.ID, &opts.AllowedApprovers)
	if err != nil {
		return err
	}
	_, _, err = d.GitLab.Projects.ChangeApprovalConfiguration(project.ID, &opts.ApprovalConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dredd) hasProjectApprovalsChanges(opts *Options, approvals *gitlab.ProjectApprovals) bool {
	changed := false
	ac := opts.ApprovalConfiguration
	aa := opts.AllowedApprovers
	if ac.ApprovalsBeforeMerge != nil {
		if *ac.ApprovalsBeforeMerge != approvals.ApprovalsBeforeMerge {
			logrus.Infof(
				"Approvals Before Merge: %d != %d",
				approvals.ApprovalsBeforeMerge,
				*ac.ApprovalsBeforeMerge,
			)
			changed = true
		}
	}
	if ac.ResetApprovalsOnPush != nil {
		if *ac.ResetApprovalsOnPush != approvals.ResetApprovalsOnPush {
			logrus.Infof(
				"Reset Approvals On Push: %v != %v",
				approvals.ResetApprovalsOnPush,
				*ac.ResetApprovalsOnPush,
			)
			changed = true
		}
	}
	if ac.DisableOverridingApproversPerMergeRequest != nil {
		if *ac.DisableOverridingApproversPerMergeRequest != approvals.DisableOverridingApproversPerMergeRequest {
			logrus.Infof(
				"Disable Overriding Approvers Per Merge Request: %v != %v",
				approvals.DisableOverridingApproversPerMergeRequest,
				*ac.DisableOverridingApproversPerMergeRequest,
			)
			changed = true
		}
	}
	if ac.MergeRequestsAuthorApproval != nil {
		if *ac.MergeRequestsAuthorApproval != approvals.MergeRequestsAuthorApproval {
			logrus.Infof(
				"Merge Requests Author Approval: %v != %v",
				approvals.MergeRequestsAuthorApproval,
				*ac.MergeRequestsAuthorApproval,
			)
			changed = true
		}
	}
	if aa.ApproverIDs != nil {
		if len(aa.ApproverIDs) != len(approvals.Approvers) {
			logrus.Infof(
				"Approvers: %d != %d",
				len(approvals.Approvers),
				len(aa.ApproverIDs),
			)
			changed = true
		}
	}
	return changed
}

func (d *Dredd) GetProjects(groupID int) ([]*gitlab.Project, error) {
	page := 1
	projects := []*gitlab.Project{}
	for {
		options := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
			IncludeSubgroups: gitlab.Bool(true),
		}
		fetchedProjects, r, err := d.GitLab.Groups.ListGroupProjects(groupID, options, nil)
		if err != nil {
			return nil, err
		}
		for _, project := range fetchedProjects {
			projects = append(projects, project)
		}
		nextPageRaw := r.Header.Get("X-Next-Page")
		if len(nextPageRaw) == 0 {
			break
		}
		nextPage, err := strconv.Atoi(nextPageRaw)
		if err != nil {
			break
		}
		page = nextPage
	}
	return projects, nil
}
