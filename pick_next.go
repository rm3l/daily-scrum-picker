package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var teamMembers = []string{
	"Armel",
	"Fortune",
	"Gennady",
	"Leanne",
	"Subhash",
	"Zbynek",
}

const stateFile = "remaining.txt"

func main() {
	rand.Seed(time.Now().UnixNano())

	remaining := loadRemaining()

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

// Load remaining names from file; if not exists, return full list shuffled
func loadRemaining() []string {
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
