package main

import (
	"fmt"
	"os"
	"strings"
)

func validateBranches(branchNames []string) {
	fmt.Println()
	fmt.Println()
	fmt.Println("Is this the correct list of branches?")
	var userInput string

	for _, branchName := range branchNames {
		fmt.Println("  ->", strings.TrimSpace(branchName))
	}

	fmt.Print("Continue? (y/n): ")
	fmt.Scan(&userInput)
	switch userInput {
	case "y":
		break
	case "n":
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid response, exiting...")
		os.Exit(0)
	}

	fmt.Println("Continuing...")
}
