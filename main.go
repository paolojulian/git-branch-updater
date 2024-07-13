package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strings"

	"paolojulian.dev/git-branch-updater/internal/git_operations"
	"paolojulian.dev/git-branch-updater/internal/logger"
	"paolojulian.dev/git-branch-updater/internal/validator"
)

const ARG_SPLITTER string = "/"

var appGitOps GitOperations = git_operations.NewGitOps()
var appLogger logger.Logger = logger.NewLogger()

func main() {
	args, err := getArgs()
	if err != nil {
		appLogger.Error(err)
	}

	appLogger.Header(1, "Fetching branches")
	if err := appGitOps.Fetch(); err != nil {
		log.Fatal(err)
	}

	appLogger.Header(2, "Convert args to full branch names")
	branchNames, err := getBranchNames(args)
	if err != nil {
		appLogger.Error(err)
	}

	validator.ValidateBranches(branchNames)

	appLogger.Header(3, "Updating branches to latest change")
	for _, branchName := range branchNames {
		pullBranch(branchName)
	}

	appLogger.Header(4, "Merge dependent branches")
	mergeDependentBranches(branchNames)

	appLogger.Header(5, "Finished")
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
	appLogger.Description("Getting all branch names (git branch -a)")

	branches, err := appGitOps.GetBranchNames()
	if err != nil {
		log.Fatal(err)
	}

	fullBranchNames := []string{}

	appLogger.Description("Mapping args to full branch names")
	for _, arg := range args {
		fullBranchName, err := getFullBranchName(arg, branches)
		if err != nil {
			appLogger.Error(err)
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
	appLogger.Description("Pulling branch: " + branchToUpdate)
	appGitOps.Switch(branchToUpdate)
	appGitOps.Pull()
}

func mergeDependentBranches(branchNames []string) {
	currentBranch := branchNames[0]
	for index, branchName := range branchNames {
		// We skip the first branch since it's the base branch
		if index == 0 {
			continue
		}
		appLogger.Description("Merging branch: " + currentBranch + " --> " + branchName)
		appGitOps.Switch(branchName)
		appGitOps.Merge(currentBranch)
		appGitOps.Push()
		currentBranch = branchName
	}
}
