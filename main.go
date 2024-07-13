package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"
)

const ARG_SPLITTER string = "/"

var gitOps GitOperations = NewGitOps()
var logger Logger = NewLogger()

func main() {
	args, err := getArgs()
	if err != nil {
		logger.Error(err)
	}

	logger.Header(1, "Fetching branches")
	if err := gitOps.Fetch(); err != nil {
		log.Fatal(err)
	}

	logger.Header(2, "Convert args to full branch names")
	branchNames, err := getBranchNames(args)
	if err != nil {
		logger.Error(err)
	}

	validateBranches(branchNames)

	logger.Header(3, "Updating branches to latest change")
	for _, branchName := range branchNames {
		pullBranch(branchName)
	}

	logger.Header(4, "Merge dependent branches")
	mergeDependentBranches(branchNames)

	logger.Header(5, "Finished")
}

func getArgs() ([]string, error) {
	argsWithoutProp := os.Args[1:]
	if len(argsWithoutProp) != 1 {
		return []string{}, errors.New("should contain exactly one arg")
	}

	args := argsWithoutProp[0]
	doesMatchFormat := regexp.MustCompile(`^([\w\d-]+)(\/[\w\d-]+)+$`).MatchString(args)
	if !doesMatchFormat {
		return []string{}, errors.New("invalid arg format, should be like 'master>developer>feature>feature-1'")
	}

	return strings.Split(args, ARG_SPLITTER), nil
}

func getBranchNames(args []string) ([]string, error) {
	logger.Description("Getting all branch names (git branch -a)")

	branches, err := gitOps.GetBranchNames()
	if err != nil {
		log.Fatal(err)
	}

	fullBranchNames := []string{}

	logger.Description("Mapping args to full branch names")
	for _, arg := range args {
		fullBranchName, err := getFullBranchName(arg, branches)
		if err != nil {
			logger.Error(err)
		}
		fullBranchNames = append(fullBranchNames, fullBranchName)
	}

	return fullBranchNames, nil
}

func getFullBranchName(shortName string, branches []string) (string, error) {
	for _, branch := range branches {
		if strings.Contains(branch, shortName) {
			trimmedSpaces := strings.TrimSpace(branch)
			removedAsterisk := strings.TrimPrefix(trimmedSpaces, "*")
			removedRemotes := strings.TrimPrefix(removedAsterisk, "remotes/")

			return strings.TrimSpace(removedRemotes), nil
		}
	}

	return "", errors.New("No branch name matches: " + shortName)
}

func pullBranch(branchName string) {
	branchToUpdate := strings.TrimPrefix(branchName, "origin/")
	logger.Description("Pulling branch: " + branchToUpdate)
	gitOps.Switch(branchToUpdate)
	gitOps.Pull()
}

func mergeDependentBranches(branchNames []string) {
	currentBranch := branchNames[0]
	for index, branchName := range branchNames {
		// We skip the first branch since it's the base branch
		if index == 0 {
			continue
		}
		logger.Description("Merging branch: " + currentBranch + " --> " + branchName)
		gitOps.Switch(branchName)
		gitOps.Merge(currentBranch)
		gitOps.Push()
		currentBranch = branchName
	}
}
