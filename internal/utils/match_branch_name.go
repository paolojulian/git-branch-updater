package match_branch_name

import (
	"regexp"
)

/**
 * MatchBranchName matches the given name with the full branch name
 *
 * e.g.
 * fullBranchName: "origin/feature/NOVA-8823/fix-bug"
 * name: "NOVA-8823"
 * returns true
 *
 * fullBranchName: "origin/feature/NOVA-8823/fix-bug"
 * name: "88"
 * returns false
 *
 * fullBranchName: "origin/feature/NOVA-8823/fix-bug"
 * name: "8823"
 * returns true
 */
func MatchBranchName(fullBranchName string, shortName string) bool {
	if shortName == "" {
		return false
	}
	// Escape the name to safely use it in regex
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(shortName) + `\b`)

	return re.MatchString(fullBranchName)
}
