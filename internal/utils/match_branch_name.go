package match_branch_name

import (
	"regexp"
	"strconv"
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

	// Check if shortName is number
	// If yes, get the latest number from the fullBranchName
	// e.g. origin/feature/NOVA-8823/partial/NOVA-8824/ui
	// should just compare NOVA-8824, not NOVA-8823
	if _, err := strconv.Atoi(shortName); err == nil {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(getLatestNumber(fullBranchName)) + `\b`)
		return re.MatchString(shortName)
	}

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
