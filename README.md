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
review-k8s: &review-k8s
  allowedApprovers:
    approverids: [1, 2, 3]
# 
options:
  devops/k8s-tools:
    <<: *global
    <<: *review-k8s
  devops/k8s-infrastructure:
    <<: *global
    <<: *review-k8s
    <<: *master-protected
  devops:
    <<: *global
    <<: *review-all
```

## TODO

* https://docs.gitlab.com/ee/api/merge_request_approvals.html#create-project-level-rule

## Links

* [GitLab Plugin system](https://docs.gitlab.com/ee/administration/plugins.html)
