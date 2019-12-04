package main

import (
	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasProjectOptionsChanges(project *gitlab.Project, opts *Options) (changed bool) {
	po := opts.ProjectOptions
	if po.ApprovalsBeforeMerge != nil {
		if *po.ApprovalsBeforeMerge != project.ApprovalsBeforeMerge {
			logrus.Infof(
				"Approvals Before Merge: %d != %d",
				*po.ApprovalsBeforeMerge,
				project.ApprovalsBeforeMerge,
			)
			changed = true
		}
	}
	if po.OnlyAllowMergeIfPipelineSucceeds != nil {
		if *po.OnlyAllowMergeIfPipelineSucceeds != project.OnlyAllowMergeIfPipelineSucceeds {
			logrus.Infof(
				"Only Allow Merge If Pipeline Succeeds: %v != %v",
				*po.OnlyAllowMergeIfPipelineSucceeds,
				project.OnlyAllowMergeIfPipelineSucceeds,
			)
			changed = true
		}
	}
	if po.OnlyAllowMergeIfAllDiscussionsAreResolved != nil {
		if *po.OnlyAllowMergeIfAllDiscussionsAreResolved != project.OnlyAllowMergeIfAllDiscussionsAreResolved {
			logrus.Infof(
				"Only Allow Merge If All Discussions Are Resolved: %v != %v",
				*po.OnlyAllowMergeIfAllDiscussionsAreResolved,
				project.OnlyAllowMergeIfAllDiscussionsAreResolved,
			)
			changed = true
		}
	}
	return changed
}

func (d *Dredd) FixProjectOptions(project *gitlab.Project, opts *Options) error {
	_, _, err := d.GitLab.Projects.EditProject(project.ID, &opts.ProjectOptions)
	return err
}
