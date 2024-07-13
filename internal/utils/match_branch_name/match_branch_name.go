package match_branch_name

import (
	"regexp"
	"strconv"
)

/**
 * Exec matches the given name with the full branch name
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
func Exec(fullBranchName string, shortName string) bool {
	if shortName == "" {
		return false
	}

	// If the arg is a number, we should get the last number in the branch name
	// e.g. feature/NOVA-8823/partial/NOVA-8824/ui
	// should match 8824 and not with 8823
	if _, err := strconv.Atoi(shortName); err == nil {
		latestNumber := getLatestNumber(fullBranchName)
		if latestNumber != "" {
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(getLatestNumber(fullBranchName)) + `\b`)
			return re.MatchString(shortName)
		}
	}

	// Arg is a string

	// Escape the name to safely use it in regex
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(shortName) + `\b`)

	return re.MatchString(fullBranchName)
}

func getLatestNumber(fullBranchName string) string {
	re := regexp.MustCompile(`\b(\d+)\b`)
	matches := re.FindAllStringSubmatch(fullBranchName, -1)
	if len(matches) == 0 {
		return ""
	}

	latestNumber := matches[len(matches)-1][1]

	return latestNumber
}
