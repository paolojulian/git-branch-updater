package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"paolojulian.dev/git-branch-updater/internal/git_operations"
	"paolojulian.dev/git-branch-updater/internal/logger"
	match_branch_name "paolojulian.dev/git-branch-updater/internal/utils"
	"paolojulian.dev/git-branch-updater/internal/validator"
)

const ARG_SPLITTER string = "/"

var APP_GIT_OPS git_operations.GitOperations = git_operations.NewGitOps()
var APP_LOGGER logger.Logger = logger.NewLogger()

var OPTIONS []string = []string{
	"--no-merge",
}
var USER_OPTIONS []string

func main() {
	args, err := getArgs()
	if err != nil {
		APP_LOGGER.Error(err)
	}

	APP_LOGGER.Header(1, "Fetching branches")
	if err := APP_GIT_OPS.Fetch(); err != nil {
		log.Fatal(err)
	}

	APP_LOGGER.Header(2, "Convert args to full branch names")
	branchNames, err := getBranchNames(args)
	if err != nil {
		APP_LOGGER.Error(err)
	}

	validator.ValidateBranches(branchNames)

	APP_LOGGER.Header(3, "Updating branches to latest change")
	for _, branchName := range branchNames {
		pullBranch(branchName)
	}

	if (USER_OPTIONS != nil) && slices.Contains(USER_OPTIONS, "--no-merge") {
		APP_LOGGER.Header(4, "Finished")
		return
	}

	APP_LOGGER.Header(4, "Merge dependent branches")
	mergeDependentBranches(branchNames)

	APP_LOGGER.Header(5, "Finished")
}

func getArgs() ([]string, error) {
	argsWithoutProp := os.Args[1:]
	if len(argsWithoutProp) == 0 {
		return []string{}, errors.New("should at least contain one arg")
	}

	args := argsWithoutProp[0]
	doesMatchFormat := regexp.MustCompile(`^([\w\d-]+)(\/[\w\d-]+)+$`).MatchString(args)
	if !doesMatchFormat {
		return []string{}, errors.New("invalid arg format, should be like 'master>developer>feature>feature-1'")
	}

	if len(argsWithoutProp) == 2 {
		option := argsWithoutProp[1]
		if (option != "") && !slices.Contains(OPTIONS, option) {
			return []string{}, errors.New("invalid option provided: " + option + ", do you mean --no-merge?")
		}
		if option != "" && slices.Contains(OPTIONS, option) {
			// The user has provided an option
			USER_OPTIONS = append(USER_OPTIONS, option)
		}
	}

	return strings.Split(args, ARG_SPLITTER), nil
}

func getBranchNames(args []string) ([]string, error) {
	APP_LOGGER.Description("Getting all branch names (git branch -a)")

	branches, err := APP_GIT_OPS.GetBranchNames()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(branches)
	fullBranchNames := []string{}

	APP_LOGGER.Description("Mapping args to full branch names")
	for _, arg := range args {
		fullBranchName, err := getFullBranchName(arg, branches)
		if err != nil {
			APP_LOGGER.Error(err)
		}
		fullBranchNames = append(fullBranchNames, fullBranchName)
	}

	return fullBranchNames, nil
}

func getFullBranchName(shortName string, branches []string) (string, error) {
	for _, branch := range branches {
		doesMatch := match_branch_name.MatchBranchName(branch, shortName)
		if doesMatch {
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
	APP_LOGGER.Description("Pulling branch: " + branchToUpdate)

	if err := APP_GIT_OPS.Switch(branchToUpdate); err != nil {
		log.Fatal(err)
	}

	if err := APP_GIT_OPS.Pull(branchToUpdate); err != nil {
		log.Fatal(err)
	}
}

func mergeDependentBranches(branchNames []string) {
	currentBranch := branchNames[0]
	for index, branchName := range branchNames {
		// We skip the first branch since it's the base branch
		if index == 0 {
			continue
		}
		APP_LOGGER.Description("Merging branch: " + currentBranch + " --> " + branchName)

		if err := APP_GIT_OPS.Switch(branchName); err != nil {
			log.Fatal(err)
		}
		if err := APP_GIT_OPS.Merge(currentBranch); err != nil {
			log.Fatal(err)
		}
		if err := APP_GIT_OPS.Push(); err != nil {
			log.Fatal(err)
		}
		currentBranch = branchName
	}
}
