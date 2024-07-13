package git_operations

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitOperations interface {
	Fetch() error
	Switch(branchName string) error
	Merge(branchName string) error
	GetBranchNames() ([]string, error)
	GetRemoteBranches() ([]string, error)
	GetCurrentBranchName() (string, error)
	Pull(branchName string) error
	Push() error
}

type GitOps struct {
}

func NewGitOps() *GitOps {
	return &GitOps{}
}

func (g *GitOps) Fetch() error {
	cmd := exec.Command("git", "fetch", "--all")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return displayGitError("failed to fetch all branches", cmd, output)
	}

	return nil
}

func (g *GitOps) Switch(branchName string) error {
	// Ensure that the branch is a local branch
	branchToSwitchTo := strings.TrimPrefix(branchName, "origin/")

	cmd := exec.Command("git", "switch", branchToSwitchTo)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return displayGitError("failed to switch to branch", cmd, output)
	}

	return nil
}

func (g *GitOps) Merge(branchName string) error {
	cmd := exec.Command("git", "merge", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return displayGitError("failed to merge branch:"+branchName, cmd, output)
	}

	return nil
}

func (g *GitOps) GetBranchNames() ([]string, error) {
	cmd := exec.Command("git", "branch", "-a")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{}, displayGitError("failed to get all branches", cmd, output)
	}

	branches := strings.Split(string(output), "\n")

	return branches, nil
}

func (g *GitOps) Pull(branchName string) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return displayGitError("failed to pull fast-forward", cmd, output)
	}

	return nil
}

func (g *GitOps) Push() error {
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return displayGitError("failed to push", cmd, output)
	}

	return nil
}

func (g *GitOps) GetRemoteBranches() ([]string, error) {
	cmd := exec.Command("git", "branch", "-r")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return []string{}, displayGitError("failed to get remote branches", cmd, output)
	}

	branches := strings.Split(string(output), "\n")
	filteredBranches := []string{}
	for _, branch := range branches {
		trimmedSpaces := strings.TrimSpace(branch)
		removedAsterisk := strings.TrimPrefix(trimmedSpaces, "*")
		removedRemotes := strings.TrimPrefix(removedAsterisk, "remotes/")

		filteredBranches = append(filteredBranches, removedRemotes)
	}

	return filteredBranches, nil
}

func (g *GitOps) GetCurrentBranchName() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", displayGitError("failed to get current branch name", cmd, output)
	}

	return strings.TrimSpace(string(output)), nil
}

func displayGitError(title string, cmd *exec.Cmd, output []byte) error {
	fmt.Println("\n******* ERROR:", title)
	fmt.Println("Command:", cmd)

	return fmt.Errorf(string(output))
}
