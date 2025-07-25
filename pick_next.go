package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// ANSI color codes - universal colors for both dark and light terminals
const (
	ColorReset = "\033[0m"

	// Standard colors that work universally
	ColorRed    = "\033[31m" // Good on both backgrounds
	ColorGreen  = "\033[32m" // Good on both backgrounds
	ColorBlue   = "\033[34m" // Good on both backgrounds
	ColorPurple = "\033[35m" // Good on both backgrounds

	// Text formatting
	Bold      = "\033[1m"
	Underline = "\033[4m"

	// Universal color combinations - tested for readability on both backgrounds
	BoldRed    = "\033[1;31m" // Excellent on both
	BoldGreen  = "\033[1;32m" // Excellent on both
	BoldBlue   = "\033[1;34m" // Excellent on both
	BoldPurple = "\033[1;35m" // Excellent on both

	// Using darker variants for even better contrast
	DarkRed   = "\033[38;5;124m" // Dark red - great on both
	DarkGreen = "\033[38;5;22m"  // Dark green - great on both
	DarkBlue  = "\033[38;5;18m"  // Dark blue - great on both
	BrightRed = "\033[91m"       // Bright red - good on both
)

const goodbyeMessage = "Goodbye!"

func getTeamFile(flagValue string) string {
	// Command-line flag takes precedence (standard practice)
	if flagValue != "" {
		return flagValue
	}
	// Environment variable as fallback
	if teamFile := os.Getenv("TEAM_FILE"); teamFile != "" {
		return teamFile
	}
	// Default fallback
	return "team.txt"
}

func getStateFile() string {
	if stateFile := os.Getenv("STATE_FILE"); stateFile != "" {
		return stateFile
	}
	// Use temporary directory for state file to ensure it's writable
	return filepath.Join(os.TempDir(), "daily-scrum-picker-remaining.txt")
}

var rootCmd = &cobra.Command{
	Use:   "daily-scrum-picker",
	Short: "A simple Go utility to fairly select the next person to speak during daily scrum/stand-up meetings",
	Run:   runApp,
}

var teamFileFlag string

func init() {
	rootCmd.Flags().StringVarP(&teamFileFlag, "team-file", "t", "", "Path to team members file, or '-' for stdin (overrides TEAM_FILE environment variable)")
}

func runApp(cmd *cobra.Command, args []string) {
	teamFile := getTeamFile(teamFileFlag)
	teamMembers, err := loadTeamMembers(teamFile)
	if err != nil {
		fmt.Printf("Error loading team members: %v\n", err)
		if teamFile != "-" {
			fmt.Printf("Please create a '%s' file with one team member name per line.\n", teamFile)
		} else {
			fmt.Printf("Please provide team member names via stdin (one per line).\n")
		}
		os.Exit(1)
	}

	if len(teamMembers) == 0 {
		if teamFile != "-" {
			fmt.Printf("No team members found in '%s'. Please add team member names (one per line).\n", teamFile)
		} else {
			fmt.Printf("No team members found in stdin. Please provide team member names (one per line).\n")
		}
		os.Exit(1)
	}

	stateFile := getStateFile()

	// Print welcome message and instructions
	fmt.Println("=== Daily Scrum Picker ===")
	if teamFile == "-" {
		fmt.Printf("Team source: stdin (%d members)\n", len(teamMembers))
	} else {
		fmt.Printf("Team file: %s (%d members)\n", teamFile, len(teamMembers))
	}
	fmt.Printf("State file: %s\n", stateFile)
	fmt.Println("\nCommands:")
	fmt.Println("  p - Pick next person")
	fmt.Println("  r - Reset and start over")
	fmt.Println("  s - Show current status")
	fmt.Println("  h - Show this help")
	fmt.Println("  q - Quit")

	// Check if we can use raw mode, otherwise fall back to buffered
	if term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println("\nPress any key (no Enter needed):")
		runRawMode(teamMembers, stateFile)
	} else {
		fmt.Println("\nType commands and press Enter:")
		runBufferedMode(teamMembers, stateFile)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runRawMode(teamMembers []string, stateFile string) {
	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Falling back to buffered mode...")
		runBufferedMode(teamMembers, stateFile)
		return
	}
	defer func() {
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			fmt.Printf("Error restoring terminal: %v\n", err)
		}
	}()

	for {
		// Print prompt and flush output
		fmt.Print("> ")

		// Read single character
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}

		char := buf[0]

		// Handle Ctrl+C
		if char == 3 {
			// Restore terminal before exiting
			if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
				fmt.Printf("Error restoring terminal: %v\n", err)
			}
			fmt.Print("\n")
			fmt.Println(goodbyeMessage)
			return
		}

		// Only process printable characters
		if char < 32 || char > 126 {
			continue
		}

		input := strings.ToLower(string(char))

		// Restore terminal temporarily for clean output
		if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
			fmt.Printf("Error restoring terminal: %v\n", err)
		}

		// Clear current line and show command
		fmt.Printf("\r> %s\n", input)

		// Handle the command
		switch input {
		case "p":
			pickNextPerson(teamMembers, stateFile)
		case "r":
			resetState(teamMembers, stateFile)
		case "s":
			showStatus(teamMembers, stateFile)
		case "h":
			showHelp()
		case "q":
			fmt.Println(goodbyeMessage)
			return
		default:
			fmt.Printf("Unknown command: '%s'. Press 'h' for help.\n", input)
		}

		fmt.Println() // Add separation

		// Re-enter raw mode for next command
		oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error re-entering raw mode, exiting...")
			return
		}
	}
}

// Fallback function for systems where raw mode doesn't work
func runBufferedMode(teamMembers []string, stateFile string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break // EOF or error
		}

		input := strings.TrimSpace(strings.ToLower(scanner.Text()))

		switch input {
		case "p", "pick":
			pickNextPerson(teamMembers, stateFile)
		case "r", "reset":
			resetState(teamMembers, stateFile)
		case "s", "status":
			showStatus(teamMembers, stateFile)
		case "h", "help":
			showHelp()
		case "q", "quit", "exit":
			fmt.Println(goodbyeMessage)
			return
		case "":
			// Empty input, just continue
			continue
		default:
			fmt.Printf("Unknown command: '%s'. Type 'h' for help.\n", input)
		}
	}
}

func pickNextPerson(teamMembers []string, stateFile string) {
	remaining := loadRemaining(teamMembers, stateFile)

	// If no one left, reset
	if len(remaining) == 0 {
		fmt.Println("Everyone has already had a turn. Resetting list...")
		remaining = shuffle(copySlice(teamMembers))
	}

	// Pick the first person (since shuffled)
	picked := remaining[0]
	remaining = remaining[1:]

	// Save updated list
	saveRemaining(remaining, stateFile)

	// Display the picked person with prominent formatting - using colors that work on both backgrounds
	fmt.Printf("🎯 Next is... %s%s%s%s\n",
		Bold, BoldBlue, picked, ColorReset)

	// Show remaining count with color that works universally
	if len(remaining) > 0 {
		fmt.Printf("%s(%d people remaining in this round)%s\n",
			BrightRed, len(remaining), ColorReset)
	} else {
		fmt.Printf("%s(That was the last person in this round)%s\n",
			BoldGreen, ColorReset)
	}
}

func resetState(teamMembers []string, stateFile string) {
	// Remove state file to reset
	if err := os.Remove(stateFile); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Warning: failed to remove state file: %v\n", err)
	}
	fmt.Printf("%s✅ State reset! All %d team members are available for selection.%s\n",
		BoldGreen, len(teamMembers), ColorReset)
}

func showStatus(teamMembers []string, stateFile string) {
	remaining := loadRemaining(teamMembers, stateFile)

	fmt.Printf("%s📊 Status:%s\n", BoldBlue, ColorReset)
	fmt.Printf("  Total team members: %s%d%s\n", DarkBlue, len(teamMembers), ColorReset)
	fmt.Printf("  Remaining this round: %s%d%s\n", BrightRed, len(remaining), ColorReset)

	if len(remaining) > 0 {
		fmt.Printf("  Still to pick: %s%s%s\n",
			DarkGreen, strings.Join(remaining, ", "), ColorReset)
	} else {
		fmt.Printf("  %sEveryone has been picked this round%s\n",
			BoldGreen, ColorReset)
	}
}

func showHelp() {
	fmt.Printf("\n%s📋 Available commands:%s\n", BoldBlue, ColorReset)
	fmt.Printf("  %sp%s, pick   - Pick the next person for daily scrum\n", BoldGreen, ColorReset)
	fmt.Printf("  %sr%s, reset  - Reset state and start over with all team members\n", BrightRed, ColorReset)
	fmt.Printf("  %ss%s, status - Show current status and remaining team members\n", BoldBlue, ColorReset)
	fmt.Printf("  %sh%s, help   - Show this help message\n", BoldPurple, ColorReset)
	fmt.Printf("  %sq%s, quit   - Exit the program\n", BoldRed, ColorReset)
	fmt.Println()
}

// Load team members from file or stdin
func loadTeamMembers(teamFile string) ([]string, error) {
	var scanner *bufio.Scanner

	if teamFile == "-" {
		// Read from stdin
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		// Read from file
		file, err := os.Open(teamFile)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Printf("Warning: failed to close file: %v\n", err)
			}
		}()
		scanner = bufio.NewScanner(file)
	}

	var members []string
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
func loadRemaining(teamMembers []string, stateFile string) []string {
	file, err := os.Open(stateFile)
	if err != nil {
		// File not found → start fresh
		return shuffle(copySlice(teamMembers))
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close file: %v\n", err)
		}
	}()

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
func saveRemaining(names []string, stateFile string) {
	if len(names) == 0 {
		// Remove file to reset
		if err := os.Remove(stateFile); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Warning: failed to remove state file: %v\n", err)
		}
		return
	}

	file, err := os.Create(stateFile)
	if err != nil {
		fmt.Printf("Error writing state file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close file: %v\n", err)
		}
	}()

	for _, name := range names {
		if _, err := file.WriteString(name + "\n"); err != nil {
			fmt.Printf("Error writing to state file: %v\n", err)
			return
		}
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
