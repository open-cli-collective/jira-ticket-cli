# jira-ticket-cli (archived)

> This repository has been consolidated into [**atlassian-cli**](https://github.com/open-cli-collective/atlassian-cli), a monorepo for all Atlassian CLI tools (Jira + Confluence).

## What happened?

`jira-ticket-cli` and `confluence-cli` were merged into a single repository — [`atlassian-cli`](https://github.com/open-cli-collective/atlassian-cli) — to share code, simplify releases, and maintain a unified development workflow. All new features and bug fixes are developed there.

## How to update

If you installed via Homebrew:

```bash
# The jira-ticket-cli cask now pulls from atlassian-cli automatically
brew upgrade jira-ticket-cli
```

If you installed via the `jtk` cask (legacy):

```bash
brew uninstall jtk
brew install open-cli-collective/tap/jira-ticket-cli
```

If you installed from GitHub releases directly, download future releases from:
https://github.com/open-cli-collective/atlassian-cli/releases

## Links

- **New repo:** [open-cli-collective/atlassian-cli](https://github.com/open-cli-collective/atlassian-cli)
- **Issues:** [atlassian-cli issues](https://github.com/open-cli-collective/atlassian-cli/issues)
