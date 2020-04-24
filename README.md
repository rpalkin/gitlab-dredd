# Dredd

Sync approvers, project options, protected branches and approval configuration in GitLab.

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
welcome-issue: &welcome-issue
  firstIssue:
    title: TODO
    description: |
      - [ ] README.md
      - [ ] CHANGELOG.md
      - [ ] Linters and tests
enable-jira: &enable-jira
  jiraIntegration:
    url: https://jira.tcsbank.ru
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

* [GitLab Plugin system](https://docs.gitlab.com/ee/administration/plugins.html)
