## Pending feature tasks for the project
- [ ] Create a CONTRIBUTING.md
- [ ] Terraform: HA setup: multiple Brightly server instances behind a load balancer all receiving updates.
- [ ] yaml files validation: warn if flag is defined in the project but not an environment
- [ ] yaml files validation: error if flag is defined in an environment but not in the project
- [ ] From Feedback: rename to 'brightly'
- [ ] From Feedback: Support multivariate string variations
- [ ] From Feedback: Support multivariate number variations

### Maybe not required for MVP:
- [ ] Terraform: validate variables (see TODOs in https://github.com/brightlyorg/terraform-aws-brightly-flags/blob/main/variables.tf)
- [ ] Terraform: Create freeform Github workflow with human input to create flags and perform other tasks ie create flag, create environment, etc.
- [ ] Use localstack to do end to end testing (if possible)
- [ ] Investigate: Change yaml parsing logic to read in maps instead of structs to better detect 0 values.

### Pending DX tasks
- [ ] temporary files and archives: keep them in memory avoiding weird filesystem bugs/flaky tests.

### Probably not required for MVP:
- [ ] Document handling of deleted environments... what happens to the flags?
- [ ] auto-publish new docker image on push to main.
- [ ] AWS credentials in GitHub actions: Use suggested approach described here: https://github.com/marketplace/actions/configure-aws-credentials-action-for-github-actions#overview
- [ ] Consider using https://pkg.go.dev/github.com/sethvargo/go-githubactions or similar
- [ ] Maybe: Publish binary artifact to Github and consume in Github Action (or publish a github action)
- [ ] Maybe: thin wrappers around SDKs that handle custom url and disabling of events (this is potentially a maintenance headache)
- [ ] At flag overview level: enable default variations (can be overridden at env level)
- [ ] Docs: (from feedback): Allow users to specify their SDKs in terraform and then only generate SDK examples based on those preferences.
