package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	gitlab "github.com/xanzy/go-gitlab"
)

type Dredd struct {
	GitLab *gitlab.Client
	Config *Config
	DryRun bool
	mode   Mode
}

func (d *Dredd) Run(mode Mode) (err error) {
	d.mode = mode
	switch mode {
	case PluginMode:
		return d.RunAsPlugin()
	case StandaloneMode:
		return d.RunAsStandalone()
	case WebhookMode:
		return d.RunAsWebhook()
	}
	return errors.New("Unsupported mode")
}

func (d *Dredd) ProcessHookPayload(r io.Reader) error {
	decoder := json.NewDecoder(r)
	hook := &gitlab.HookEvent{}
	err := decoder.Decode(&hook)
	if err != nil {
		return err
	}
	logrus.Debugf("Payload received: %#v", hook)
	if hook.EventName != "project_create" {
		logrus.Debugf("Event %s was skipped", hook.EventName)
		return nil
	}
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

func (d *Dredd) RunAsPlugin() error {
	reader := bufio.NewReader(os.Stdin)
	err := d.ProcessHookPayload(reader)
	if err != nil {
		return err
	}
	return nil
}

func (d *Dredd) RunAsWebhook() error {
	http.HandleFunc("/dredd", func(w http.ResponseWriter, r *http.Request) {
		err := d.ProcessHookPayload(r.Body)
		if err != nil {
			logrus.Errorf("Error processing hook payload: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := w.Write([]byte(`ok`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.Handle("/metrics", promhttp.Handler())
	logrus.Infof("Listen address: %s", d.Config.ListenAddress)
	return http.ListenAndServe(d.Config.ListenAddress, nil)
}

func (d *Dredd) RunAsStandalone() error {
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
	if d.HasProjectApprovalRuleChanges(project, opts) && !d.DryRun {
		err := d.FixProjectApprovalRule(project, opts)
		if err != nil {
			return err
		}
	}
	if d.mode != StandaloneMode {
		if d.HasFirstIssueChanges(project, opts) && !d.DryRun {
			err := d.CreateFirstIssue(project, opts)
			if err != nil {
				return err
			}
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
		projects = append(projects, fetchedProjects...)
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
