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
	if err == nil {
		return nil
	}

	hasRemoteBranch, err := checkIfHasRemoteBranch(branchName)
	println("Has remote branch:", hasRemoteBranch)
	if err != nil {
		return displayGitError("failed to pull fast-forward", cmd, output)
	}

	return displayGitError("failed to pull fast-forward", cmd, output)
}

func checkIfHasRemoteBranch(branchName string) (bool, error) {
	// git ls-remote --heads origin refs/heads/[branch-name]
	cmd := exec.Command("git", "ls-remote", "--heads", "origin", "refs/heads/"+branchName)
	output, err := cmd.CombinedOutput()
	println("checkIfHasRemoteBranch output", string(output))
	println("checkIfHasRemoteBranch error", err.Error())
	if err != nil {
		return false, err
	}

	println("Output", string(output))
	return true, nil
}

func (g *GitOps) Push() error {
	cmd := exec.Command("git", "push", "-u", "origin", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return displayGitError("failed to push", cmd, output)
	}

	return nil
}

func displayGitError(title string, cmd *exec.Cmd, output []byte) error {
	fmt.Println("\n******* ERROR:", title)
	fmt.Println("Command:", cmd)

	return fmt.Errorf(string(output))
}
