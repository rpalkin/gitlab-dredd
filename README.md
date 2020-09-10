# Dredd

[![Build Status](https://travis-ci.com/leominov/gitlab-dredd.svg?branch=master)](https://travis-ci.com/leominov/gitlab-dredd)

Supported Gitlab project settings:
 - Sync approvers
 - Project options
 - Protected branches
 - Approval configuration
 - Jira integration
 - Push rules
 
 Yaml configuration supports environment variables in ${VAR} notation

## Configuration example

```yaml
---
listenAddress: 0.0.0.0:8080
gitlabEndpoint: https://gitlab.local
gitlabToken: ABCD
groupIDs: [1, 2, 3]
# 
global: &global
  projectOptions:
    onlyallowmergeifpipelinesucceeds: true
    onlyallowmergeifalldiscussionsareresolved: true
    approvalsbeforemerge: 1
  approvalConfiguration:
    approvalsbeforemerge: 1
    resetapprovalsonpush: true
    disableoverridingapproverspermergerequest: false
    mergerequestsauthorapproval: false
  allowedApprovers:
    approvergroupids: []
master-protected: &master-protected
  repositoryBranches:
    - name: master
      pushaccesslevel: 0
      mergeaccesslevel: 30
review-all: &review-all
  allowedApprovers:
    approverids: [1, 2, 3, 4, 5]
  approvalRule:
    name: Default
    approvalsrequired: 1
    userids: [1, 2, 3, 4, 5]
review-k8s: &review-k8s
  allowedApprovers:
    approverids: [1, 2, 3]
  approvalRule:
    name: Default
    approvalsrequired: 1
    userids: [1, 2, 3]
branch-and-commit-with-issue: &branch-and-commit-with-issue
  pushRules:
    branchnameregex: ^DEV-[\d]+.*
    commitmessageregex: ^DEV-[\d]+.*
welcome-issue: &welcome-issue
  firstIssue:
    title: TODO
    description: |
      - [ ] README.md
      - [ ] CHANGELOG.md
      - [ ] Linters and tests
enable-jira: &enable-jira
  jiraIntegration:
    url: https://jira.example.org
    username: USER
    password: PASSWORD
    active: true
    commitevents: false
    mergerequestsevents: true
    commentoneventenabled: false
# 
options:
  devops/k8s-tools:
    <<: *global
    <<: *review-k8s
    <<: *welcome-issue
    <<: *enable-jira
    <<: *branch-and-commit-with-issue
  devops/k8s-infrastructure:
    <<: *global
    <<: *review-k8s
    <<: *master-protected
    <<: *enable-jira
  devops:
    <<: *global
    <<: *review-all
```

## Links

* [Docker Hub](https://hub.docker.com/repository/docker/leominov/gitlab-dredd)
* [GitLab Plugin system](https://docs.gitlab.com/ee/administration/plugins.html)
