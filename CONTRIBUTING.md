# Contributing Guidelines

Thank you for your interest in contributing to Kaia. As an open source project, Kaia is always open to the developer community and we welcome your contribution. Please read the guideline below and follow it in all interactions with the project.

## How to Contribute

1. Read this [contributing document](./CONTRIBUTING.md).
2. Sign [Contributor Licensing Agreement (CLA)](#contributor-license-agreement-cla) in the **PR comment**.
3. Before making a PR, please make sure to check the following:
  - You should format your change. Kaia uses a stricter go format [golangci-lint](https://github.com/golangci/golangci-lint). Please take a look [how to lint your change](HOW-TO-LINT-YOUR-CHANGE.md).
  - You should make sure you fully tested the code via running `make test`.
  - You should make sure the PR targets the `dev` branch.
  - You should make sure the PR is not too large. Your PR may be rejected if the changed LOC is over 1,000. It is recommended to split the PR into smaller ones.
4. After submitting the PR, wait for code review and approval. The reviewer may ask you for additional commits or changes.
5. Once the change has been approved, the PR will be merged by the project moderator.
6. After merging the PR, we close the pull request. You can then delete the now obsolete branch.

## Types of Contribution

There are various ways to contribute and participate. Please read the guidelines below regarding the process of each type of contribution.

- [Issues and Bugs](#issues-and-bugs)
- [Feature Requests](#feature-requests)
- [Code Contribution](#code-contribution)

### Issues and Bugs

If you find a bug or other issues in Kaia, please [submit an issue](https://github.com/kaiachain/kaia/issues). If the bug is related to security, please follow [SECURITY.md](./SECURITY.md). Before submitting an issue, please invest some extra time to figure out that:

- The issue is not a duplicate issue.
- The issue has not been fixed in the latest release of Kaia.

Please do not use the issue tracker for personal support requests. Use [Kaia Dev Forum](https://devforum.kaia.io/) for the personal support requests.

When you report a bug, please make sure that your report has the following information.

- Steps to reproduce the issue.
- A clear and complete description of the issue.
- Code and/or screen captures are highly recommended.

After confirming your report meets the above criteria, [submit the issue](https://github.com/kaiachain/kaia/issues).

### Feature Requests

You can also use the [issue tracker](https://github.com/kaiachain/kaia/issues) to request a new feature or enhancement. Note that any code contribution without an issue link will not be accepted. Please submit an issue explaining your proposal first so that the Kaia community can fully understand and discuss the idea.

### Code Contribution

Please follow the coding style and quality requirements to satisfy the product standards. You must follow the coding style as best as you can when submitting code. Take note of naming conventions, separation of concerns, and formatting rules.

The go implementation of Kaia uses [godoc](https://pkg.go.dev/golang.org/x/tools/cmd/godoc)
to document its source code. For the guideline of official Go language, please
refer to the following websites:
- https://go.dev/doc/effective_go#commentary
- https://go.dev/blog/godoc

## Versioning Policy

Kaia follows [Semantic Versioning](https://semver.org/) format `v{ MAJOR }.{ MINOR }.{ PATCH }`. Increment the:

- **MAJOR** version when both conditions are met: (1) when a breaking change (hard fork) occurs, and (2) when tokenomics/governance is affected.
- **MINOR** version for most regular client updates.
- **PATCH** version for simple/urgent bug fixes, improvements, or hard fork activation block number updates.

## Contributor License Agreement (CLA)

Keep in mind when you submit your pull request, you will need to sign the [CLA](https://gist.github.com/kaiachain-dev/bbf65cc330275c057463c4c94ce787a6) via the PR comment for legal purposes. You will have to sign the CLA just one time, either as an individual or corporation.

You will be prompted to sign the agreement by CLA Assistant (bot) when you open a Pull Request for the first time.
