package check_if_has_remote_branch

import "testing"

func TestCheckIfHasRemoteBranch(t *testing.T) {
	branches := []string{
		"origin/feature/NOVA-8823/fix-bug",
		"origin/feature/NOVA-8823/partial/NOVA-8824/ui",
		"origin/feature/NOVA-8823/partial/NOVA-8824/ui",
		"origin/feature/NOVA-8823/fix-bug",
	}
	branchName := "feature/NOVA-8823/fix-bug"

	expected := true
	result := Exec(branches, branchName)

	if !result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}

func Test_if_remote_branch_does_not_exist(t *testing.T) {
	branches := []string{
		"origin/feature/NOVA-8823/fix-bug",
		"origin/feature/NOVA-8823/partial/NOVA-8824/ui",
		"origin/feature/NOVA-8823/partial/NOVA-8824/ui",
		"origin/feature/NOVA-8823/fix-bug",
	}
	branchName := "feature/NOVA-123213/should-not-exist"

	expected := false
	result := Exec(branches, branchName)

	if result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}
