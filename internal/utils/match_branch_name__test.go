package match_branch_name

import (
	"testing"
)

func TestMatchBranchName(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/fix-bug"
	shortName := "NOVA-8823"

	result := MatchBranchName(fullBranchName, shortName)

	if !result {
		t.Errorf("Expected true, got %t", result)
	}
}

func TestMatchBranchName2(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/fix-bug"
	shortName := "8823"

	result := MatchBranchName(fullBranchName, shortName)

	if !result {
		t.Errorf("Expected true, got %t", result)
	}
}

func TestMatchBranchName3(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/partial/NOVA-8824/ui"
	shortName := "8824"

	result := MatchBranchName(fullBranchName, shortName)

	if !result {
		t.Errorf("Expected true, got %t", result)
	}
}

func TestMatchBranchName4(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/partial/NOVA-8824/ui"
	shortName := "NOVA-8824"

	result := MatchBranchName(fullBranchName, shortName)

	if !result {
		t.Errorf("Expected true, got %t", result)
	}
}

func TestNoMatchBranchName(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/fix-bug"
	shortName := "88"

	result := MatchBranchName(fullBranchName, shortName)

	if result {
		t.Errorf("Expected false, got %t", result)
	}
}

func TestNoMatchBranchName2(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/fix-bug"
	shortName := "NOVA-882"

	result := MatchBranchName(fullBranchName, shortName)

	if result {
		t.Errorf("Expected false, got %t", result)
	}
}

func TestNoMatchBranchName3(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/fix-bug"
	shortName := "NOVA-8824"

	result := MatchBranchName(fullBranchName, shortName)

	if result {
		t.Errorf("Expected false, got %t", result)
	}
}

func TestWithWhiteSpaces(t *testing.T) {
	fullBranchName := " origin/feature/NOVA-8823/fix-bug "
	shortName := "NOVA-8823"

	result := MatchBranchName(fullBranchName, shortName)

	if !result {
		t.Errorf("Expected true, got %t", result)
	}
}

func TestWithSubTas(t *testing.T) {
	fullBranchName := "origin/feature/NOVA-8823/partial/NOVA-8824/ui"
	shortName := "8823"

	const expected bool = false
	result := MatchBranchName(fullBranchName, shortName)

	if result != expected {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}

func TestStringToNumber(t *testing.T) {
	fullBranchName := "dev"
	shortName := "8823"

	expected := false
	result := MatchBranchName(fullBranchName, shortName)

	if expected != result {
		t.Errorf("Expected %t, got %t", expected, result)
	}
}
