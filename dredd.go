package main

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
)

type Dredd struct {
	GitLab *gitlab.Client
	Config *Config
	DryRun bool
}

func (d *Dredd) Hook() error {
	hook, err := GetStdinHookPayload()
	if err != nil {
		return fmt.Errorf("Failed to process payload: %v", err)
	}
	if hook.EventName != "project_create" {
		return nil
	}
	logrus.Debugf("Payload received: %#v", hook)
	project, _, err := d.GitLab.Projects.GetProject(hook.ProjectID, nil)
	if err != nil {
		return err
	}
	err = d.processProjects([]*gitlab.Project{project})
	if err != nil {
		return err
	}
	return nil
}

func (d *Dredd) Run() error {
	logrus.Info("Requesting projects list...")
	var projects []*gitlab.Project
	for _, group := range d.Config.GroupIDs {
		ps, err := d.GetProjects(group)
		if err != nil {
			return err
		}
		projects = append(projects, ps...)
	}
	err := d.processProjects(projects)
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
	logrus.Infof("Processing %s project...", project.PathWithNamespace)
	opts := d.Config.OptionsByPath(project.PathWithNamespace)
	if opts == nil {
		return nil
	}
	if d.HasProjectOptionsChanges(project, opts) && !d.DryRun {
		err := d.FixProjectOptions(project, opts)
		if err != nil {
			return err
		}
	}
	if d.HasAllowedApproversChanges(project, opts) && !d.DryRun {
		err := d.FixAllowedApprovers(project, opts)
		if err != nil {
			return err
		}
	}
	if d.HasApprovalConfigurationChanges(project, opts) && !d.DryRun {
		err := d.FixApprovalConfiguration(project, opts)
		if err != nil {
			return err
		}
	}
	if d.HasRepositoryBranchesOptionsChanges(project, opts) && !d.DryRun {
		err := d.FixRepositoryBranchesOptions(project, opts)
		if err != nil {
			return err
		}
	}
	return nil
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
