package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const stateFile = "remaining.txt"

func getTeamFile() string {
	if teamFile := os.Getenv("TEAM_FILE"); teamFile != "" {
		return teamFile
	}
	return "team.txt"
}

func main() {
	rand.Seed(time.Now().UnixNano())

	teamFile := getTeamFile()
	teamMembers, err := loadTeamMembers(teamFile)
	if err != nil {
		fmt.Printf("Error loading team members: %v\n", err)
		fmt.Printf("Please create a '%s' file with one team member name per line.\n", teamFile)
		os.Exit(1)
	}

	if len(teamMembers) == 0 {
		fmt.Printf("No team members found in '%s'. Please add team member names (one per line).\n", teamFile)
		os.Exit(1)
	}

	remaining := loadRemaining(teamMembers)

	// If no one left, reset
	if len(remaining) == 0 {
		fmt.Println("Everyone has already had a turn. Resetting list...")
		remaining = shuffle(copySlice(teamMembers))
	}

	// Pick the first person (since shuffled)
	picked := remaining[0]
	remaining = remaining[1:]

	// Save updated list
	saveRemaining(remaining)

	fmt.Printf("Next is... %s\n", picked)
}

// Load team members from file
func loadTeamMembers(teamFile string) ([]string, error) {
	file, err := os.Open(teamFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var members []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" && !strings.HasPrefix(name, "#") {
			members = append(members, name)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// Load remaining names from file; if not exists, return full list shuffled
func loadRemaining(teamMembers []string) []string {
	file, err := os.Open(stateFile)
	if err != nil {
		// File not found → start fresh
		return shuffle(copySlice(teamMembers))
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			names = append(names, name)
		}
	}
	return names
}

// Save remaining names to file
func saveRemaining(names []string) {
	if len(names) == 0 {
		// Remove file to reset
		os.Remove(stateFile)
		return
	}

	file, err := os.Create(stateFile)
	if err != nil {
		fmt.Printf("Error writing state file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	for _, name := range names {
		file.WriteString(name + "\n")
	}
}

// Shuffle a slice
func shuffle(slice []string) []string {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
	return slice
}

// Helper to copy slice
func copySlice(slice []string) []string {
	newSlice := make([]string, len(slice))
	copy(newSlice, slice)
	return newSlice
}
