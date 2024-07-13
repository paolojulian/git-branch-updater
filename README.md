# git-branch-updater

TLDR: A Go-based CLI tool that simplifies updating dependent Git branches. Ideal for developers managing complex branching strategies, it significantly reduces the time needed to update feature branches.

## Key Features

- **Automatic Branch Detection**: Easily identify and manage dependent branches.
- **Cascading Updates**: Propagate changes through a chain of dependent branches.
- **Smart Merging**: Intelligently merge changes while minimizing conflicts.
- **User-Friendly Interface**: Simple CLI commands for ease of use.
- **Flexible Branch Naming**: Support for various branch naming conventions.
- **Error Handling**: Robust error detection and reporting.
- **Fast and Efficient**: Optimized performance for quick updates.

## Use Cases

Given multiple dependent branches:

```
main
-> dev
-> feat/81/task-list
-> feat/82/ui
-> feat/83/integrate-api
-> feat/84/tests

```
Updating `feat/84/tests` from the latest `main` can be tedious. With `git-branch-updater`, it's a single command:


```
$ ./git-branch-updater main/dev/81/82/83/84
```

### Result
```
$ ./git-branch-updater main/dev/81/82/83/84

-- 1. Fetching branches
-- 2. Convert args to full branch names
---- Getting all branch names (git branch -a)
---- Mapping args to full branch names

Is this the correct list of branches?
  -> main
  -> dev
  -> feat/81/task-list
  -> feat/82/ui
  -> feat/83/integrate-api
  -> feat/84/tests
Continue? (y/n): y 
Continuing...

-- 3. Updating branches to latest change
---- Pulling branch: main
---- Pulling branch: dev
---- Pulling branch: feat/81/task-list
---- Pulling branch: feat/82/ui
---- Pulling branch: feat/83/integrate-api
---- Pulling branch: feat/84/tests
-- 4. Merge dependent branches
---- Merging branch: main --> dev
---- Merging branch: dev --> feat/81/task-list
---- Merging branch: feat/81/task-list --> feat/82/ui
---- Merging branch: feat/82/ui --> feat/83/integrate-api
---- Merging branch: feat/83/integrate-api --> feat/84/tests
-- 5. Finished%                     
```

## Options
- `--no-merge` : This only pulls the latest changes into its own feature branch, this will not merge dependent branch

Example: If you just want to update your local branches (`main`, `staging`, `dev`)
```
$ git-branch-updater main/staging/dev --no-merge
```