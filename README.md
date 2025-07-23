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

## Customization

To modify the team members, edit the `teamMembers` slice in `pick_next.go`:

```go
var teamMembers = []string{
    "Your",
    "Team",
    "Members",
    "Here",
}
```

Then rebuild the application:
```bash
go build -o scrum-picker pick_next.go
```

## Technical Details

- **Language**: Go 1.24.4
- **Dependencies**: Standard library only
- **State Storage**: Plain text file (`remaining.txt`)
- **Randomization**: Uses `math/rand` with time-based seeding

## File Structure

```
daily-scrum-picker/
├── go.mod              # Go module definition
├── pick_next.go        # Main application code
├── remaining.txt       # State file (created automatically)
└── README.md          # This file
``` 