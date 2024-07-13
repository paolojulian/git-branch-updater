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
	"paolojulian.dev/git-branch-updater/internal/utils/check_if_has_remote_branch"
	"paolojulian.dev/git-branch-updater/internal/utils/match_branch_name"
	"paolojulian.dev/git-branch-updater/internal/validator"
)

const ARG_SPLITTER string = "/"

var APP_GIT_OPS git_operations.GitOperations = git_operations.NewGitOps()
var APP_LOGGER logger.Logger = logger.NewLogger()

var OPTIONS []string = []string{
	"--update-only",
}
var USER_OPTIONS []string
var CURRENT_BRANCH string

func main() {
	args, err := getArgs()
	if err != nil {
		APP_LOGGER.Error(err)
	}

	currentBranch, err := APP_GIT_OPS.GetCurrentBranchName()
	CURRENT_BRANCH = currentBranch
	if err != nil {
		log.Fatal(err)
	}

	APP_LOGGER.Header("Fetching branches")
	if err := APP_GIT_OPS.Fetch(); err != nil {
		log.Fatal(err)
	}

	APP_LOGGER.Header("Convert args to full branch names")
	branchNames, err := getBranchNames(args)
	if err != nil {
		APP_LOGGER.Error(err)
	}

	validator.ValidateBranches(branchNames)

	APP_LOGGER.Header("Updating branches to latest change")
	remoteBranches, err := APP_GIT_OPS.GetRemoteBranches()
	if err != nil {
		log.Fatal(err)
	}

	for _, branchName := range branchNames {
		hasRemoteBranch := check_if_has_remote_branch.Exec(remoteBranches, branchName)
		pullBranch(branchName, hasRemoteBranch)
	}

	if (USER_OPTIONS != nil) && slices.Contains(USER_OPTIONS, "--update-only") {
		switchToCurrentBranch()
		APP_LOGGER.Header("Finished")
		return
	}

	APP_LOGGER.Header("Merge dependent branches")
	mergeDependentBranches(branchNames)

	switchToCurrentBranch()
	APP_LOGGER.Header("Finished")
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
			return []string{}, errors.New("invalid option provided: " + option + ", do you mean --update-only?")
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
		doesMatch := match_branch_name.Exec(branch, shortName)
		if doesMatch {
			trimmedSpaces := strings.TrimSpace(branch)
			removedAsterisk := strings.TrimPrefix(trimmedSpaces, "*")
			removedRemotes := strings.TrimPrefix(removedAsterisk, "remotes/")

			return strings.TrimSpace(removedRemotes), nil
		}
	}

	return "", errors.New("No branch name matches: " + shortName)
}

func pullBranch(branchName string, hasRemoteBranch bool) {
	branchToUpdate := strings.TrimPrefix(branchName, "origin/")
	APP_LOGGER.Description("Pulling branch: " + branchToUpdate)

	if !hasRemoteBranch {
		APP_LOGGER.Description("Branch not found in remote, skipping")
		return
	}

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

func switchToCurrentBranch() {
	if CURRENT_BRANCH == "" {
		return
	}

	err := APP_GIT_OPS.Switch(CURRENT_BRANCH)
	if err != nil {
		log.Fatal(err)
	}
}
