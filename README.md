# git-branch-updater

TLDR; `git-branch-updater` is a simple Command Line Interface (CLI) tool written in Go, designed to update dependent feature branches. This tool is perfect for developers working on projects with complex branching strategies, helping to quicken the time updating one feature branch to the other feature branch

## Key Features

- **Automatic Branch Detection**: Easily identify and manage dependent branches.
- **Cascading Updates**: Propagate changes through a chain of dependent branches.
- **Smart Merging**: Intelligently merge changes while minimizing conflicts.
- **User-Friendly Interface**: Simple CLI commands for ease of use.
- **Flexible Branch Naming**: Support for various branch naming conventions.
- **Error Handling**: Robust error detection and reporting.
- **Fast and Efficient**: Optimized performance for quick updates.

## Use Cases

Given you have multiple dependent branches

### Example
We need to update `feat/84/tests` from latest `main`, but there are too many dependent branch.

```
main
-> dev
-> feat/81/task-list
-> feat/82/ui
-> feat/83/integrate-api
-> feat/84/tests
```

It is quite painful to switch to them one-by-one, with `git-branch-updater` you can do it only on one command.

### Using git-branch-updater
```
$ ./git-branch-updater main/dev/81/82/83/84
```

### Result
```
```


## Getting Started

Check out our [Installation Guide](link-to-installation-guide) and [Usage Examples](link-to-usage-examples) to start optimizing your Git branch management today!
