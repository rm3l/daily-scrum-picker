# Daily Scrum Picker

A simple Go utility to fairly select the next person to speak during daily scrum/stand-up meetings.

## Overview

This tool ensures fair rotation of team members during daily stand-ups by maintaining state between runs and automatically resetting when everyone has had a turn.

## How It Works

1. **Fair Rotation**: Tracks who hasn't spoken yet using a state file (`remaining.txt`)
2. **Random Shuffling**: When everyone has had a turn, shuffles the team list for the next cycle
3. **Persistent State**: Remembers selections between runs so you can use it daily
4. **Automatic Reset**: When the list is empty, automatically starts a new randomized cycle

## Installation

### Prerequisites
- Go 1.24.4 or later

### Setup
```bash
git clone https://gitlab.cee.redhat.com/asoro/rhdh-install-daily-scrum-picker.git
cd daily-scrum-picker
```

## Usage

Simply run with Go:
```bash
go run pick_next.go
```

Alternatively, you can build and run the executable:
```bash
go build -o scrum-picker pick_next.go
./scrum-picker
```

### Output Examples

**First run or after reset:**
```
Everyone has already had a turn. Resetting list...
Next is... Gennady
```

**Regular run:**
```
Next is... Leanne
```

## State Management

- **State File**: `remaining.txt` - stores names of team members who haven't been picked yet
- **Auto-cleanup**: The state file is automatically removed when empty and recreated as needed
- **Reset Logic**: When no one is left in the remaining list, the tool automatically shuffles and restarts

## Configuration

### Team Members

Team members are configured via a team file. By default, the tool looks for `team.txt`, but you can specify a different file using the `TEAM_FILE` environment variable.

Create or edit the team file with one team member name per line:

```
# Daily Scrum Team Members
# Add one team member name per line
# Lines starting with # are ignored

Alice
Bob
Charlie
Diana
```

#### Using a Custom Team File

You can specify a different team file location using the `TEAM_FILE` environment variable:

```bash
# Use a different file
export TEAM_FILE="my-team.txt"
go run pick_next.go

# Or inline
TEAM_FILE="teams/backend-team.txt" go run pick_next.go
```

The tool will automatically read from the specified file when started. If the file doesn't exist, you'll get a helpful error message.

## Technical Details

- **Language**: Go 1.24.4
- **Dependencies**: Standard library only
- **Configuration**: Team file (default: `team.txt`, configurable via `TEAM_FILE` env var)
- **State Storage**: Plain text file (`remaining.txt`)
- **Randomization**: Uses `math/rand` with time-based seeding

## File Structure

```
daily-scrum-picker/
├── go.mod              # Go module definition
├── pick_next.go        # Main application code
├── team.txt            # Team members configuration
├── remaining.txt       # State file (created automatically)
└── README.md          # This file
``` 