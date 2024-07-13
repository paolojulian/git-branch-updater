package check_if_has_remote_branch

import "strings"

func Exec(remoteBranchNames []string, localBranchName string) bool {
	for _, remoteBranchName := range remoteBranchNames {
		if remoteBranchName == localBranchName || strings.HasSuffix(remoteBranchName, "/"+localBranchName) {
			return true
		}
	}
	return false
}
