# ZGit

ZGit is a Git workflow enhancement tool that automates common tasks and improves developer productivity. It provides automatic ticket tracking, template-based commit messages, and advanced branch synchronization features.

## Features

- **Automatic ticket extraction** - Extracts ticket numbers from branch names using configurable regex patterns
- **Template-based commit messages** - Automatically formats commit messages with ticket prefixes
- **Force-pull support** - Safely syncs local branches after force-pushed remote branches
- **Repository-specific configuration** - Supports both global and per-repository settings
- **Flexible branch patterns** - Configure custom regex patterns to match your team's naming conventions
- **Git command pass-through** - Any unknown command is automatically passed to git, making zgit a drop-in replacement

## Install or Upgrade

If you are root user, you can run the following command:

```bash
curl -sL https://gundamz.net/zgit/install.sh | bash
```

If you are not root user, you can run the following command:

```bash
curl -sL https://gundamz.net/zgit/install.sh | sudo bash
```

## Uninstall

```bash
rm -rf /usr/local/bin/zgit
rm -rf ~/.config/zgit
```

## Configuration

ZGit looks for `config.yaml` in the current directory or `~/.config/zgit/config.yaml`.

### Example Configuration

```yaml
global:
  branches:
    - usr/[^/]+/(?P<ticket>JIRA-\d+)
  commit:
    message: "[{{.Ticket}}] {{.Message}}"
repos:
  - name: owner/repo
    branches:
      - usr/[^/]+/(?P<ticket>PROJ-\d+)
```

### Configuration Options

- **global.branches** - Array of regex patterns to match branch names and extract ticket numbers
- **global.commit.message** - Template for commit messages using `{{.Ticket}}` and `{{.Message}}` placeholders
- **repos** - Array of repository-specific configurations that override global settings

## Usage

### Commit with Automatic Ticket Prefix

The `commit` command extracts the ticket number from your current branch name and automatically formats the commit message when using the `-m` flag.

**With `-m` flag (automatic ticket prefix):**

```bash
zgit commit -m "fix login bug"
```

If you're on branch `usr/john/JIRA-1234`, this will execute:

```bash
git commit -m "[JIRA-1234] fix login bug"
```

**Without `-m` flag (direct git pass-through):**

```bash
zgit commit --amend
# Opens editor, no ticket processing
```

**With additional git flags:**

```bash
zgit commit -m "fix bug" --no-verify
# Results in: git commit -m "[JIRA-1234] fix bug" --no-verify
```

### Force Pull

The `force-pull` command safely syncs your local branch with a force-pushed remote branch by recreating the local branch from origin.

```bash
zgit force-pull
```

### Using Any Git Command

ZGit acts as a transparent wrapper for git. Any command not explicitly handled by zgit (like `commit`, `force-pull`, `init`, `version`) is automatically passed to git:

```bash
zgit status
zgit log --oneline -n 10
zgit push origin main
zgit checkout -b feature/new-feature
zgit rebase -i HEAD~3
```

This means you can use `zgit` as a complete replacement for `git` in your workflow.

### Working with Different Repositories

You can specify a repository directory using the `-C` flag:

```bash
zgit -C /path/to/repo commit -m "update documentation"
zgit -C /path/to/repo status
```

## Example Workflows

1. **Standard commit with ticket tracking**

   ```bash
   # On branch usr/alice/JIRA-5678
   zgit add .
   zgit commit -m "implement user authentication"
   # Results in: git commit -m "[JIRA-5678] implement user authentication"
   ```

2. **Commit with git flags**

   ```bash
   # On branch usr/alice/JIRA-5678
   zgit commit -m "fix typo" --no-verify
   # Results in: git commit -m "[JIRA-5678] fix typo" --no-verify
   ```

3. **Amend without ticket processing**

   ```bash
   zgit commit --amend
   # Opens editor, passes directly to git
   ```

4. **Sync after force push**

   ```bash
   # Remote branch was force-pushed
   zgit force-pull
   # Safely recreates local branch from origin
   ```

5. **Complete git replacement workflow**

   ```bash
   zgit status
   zgit add .
   zgit commit -m "add new feature"
   zgit push origin main
   zgit log --oneline -n 5
   ```

6. **Multi-repository workflow**

   ```bash
   zgit -C ~/projects/frontend commit -m "update UI"
   zgit -C ~/projects/backend commit -m "add API endpoint"
   ```

## Reference

* [zgit](https://github.com/zhaojunlucky/zgit)
