package main

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

func (d *Dredd) HasRepositoryBranchesOptionsChanges(project *gitlab.Project, opts *Options) (changed bool) {
	logrus.Info("Branch options is constantly refreshed")
	return true
}

func (d *Dredd) FixRepositoryBranchesOptions(project *gitlab.Project, opts *Options) error {
	for _, branchOpt := range opts.RepositoryBranches {
	reprotect:
		logrus.Infof("Protect branch: %s", *branchOpt.Name)
		_, _, err := d.GitLab.ProtectedBranches.ProtectRepositoryBranches(project.ID, branchOpt)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				logrus.Infof("Unprotect branch: %s", *branchOpt.Name)
				d.GitLab.ProtectedBranches.UnprotectRepositoryBranches(project.ID, *branchOpt.Name)
				goto reprotect
			}
			return err
		}
	}
	return nil
}
