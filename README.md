# Daily Scrum Picker

A simple Go utility to fairly select the next person to speak during daily scrum/stand-up meetings.

## Overview

This tool ensures fair rotation of team members during daily stand-ups by automatically tracking who has spoken and resetting when everyone has had a turn.

## How It Works

1. **Fair Rotation**: Tracks who hasn't spoken yet automatically
2. **Random Shuffling**: When everyone has had a turn, shuffles the team list for the next cycle
3. **Persistent Tracking**: Remembers selections between runs so you can use it daily
4. **Automatic Reset**: When the list is empty, automatically starts a new randomized cycle

## Installation

### Prerequisites
- Go 1.24+

### Setup
```bash
git clone https://gitlab.cee.redhat.com/asoro/rhdh-install-daily-scrum-picker.git && cd rhdh-install-daily-scrum-picker
```

## Usage

### Local Development

Simply run with Go:
```bash
go run pick_next.go
```

Alternatively, you can build and run the executable:
```bash
go build -o scrum-picker pick_next.go
./scrum-picker
```

### Container Usage

The application is also available as a container image on Quay.io. This allows you to use the tool without cloning the repository or installing Go.

#### Interactive Mode

The tool runs in interactive mode with **single-keypress commands** - no need to press Enter:

```bash
# Run with the default team.txt file included in the container
podman run -it --rm quay.io/asoro/rhdh-install-daily-scrum-picker:latest
```

**Available commands (single keypress):**
- **`p`** - Pick the next person for daily scrum
- **`r`** - Reset and start over with all team members  
- **`s`** - Show current status and remaining team members
- **`h`** - Show help message
- **`q`** - Exit the program

**Note:** 
- Use the `-it` flags to enable interactive mode with proper terminal support
- Commands respond immediately without pressing Enter
- Fallback to Enter-required mode if raw terminal access is unavailable

#### Using a Custom Team File

```bash
# Mount your custom team file
podman run -it --rm -v ./my-team.txt:/app/team.txt quay.io/asoro/rhdh-install-daily-scrum-picker:latest

# Or use a different file path with environment variable
podman run -it --rm -v ./teams:/app/teams -e TEAM_FILE=teams/backend.txt quay.io/asoro/rhdh-install-daily-scrum-picker:latest
```

### Output Examples

**Interactive session:**
```
=== Daily Scrum Picker ===
Team file: team.txt (6 members)

Commands:
  p - Pick next person
  r - Reset and start over
  s - Show current status
  h - Show this help
  q - Quit

Press any key (no Enter needed):
> p
🎯 Next is... Alice
(5 people remaining in this round)
> p
🎯 Next is... Charlie
(4 people remaining in this round)
> s
📊 Status:
  Total team members: 6
  Remaining this round: 4
  Still to pick: Bob, Diana, Frank, Grace
> q
Goodbye!
```

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

#### Environment Variables

The tool supports the `TEAM_FILE` environment variable for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `TEAM_FILE` | Path to team members file | `team.txt` |

**Examples:**

```bash
# Use a different team file
export TEAM_FILE="my-team.txt"
go run pick_next.go

# Use specific team file
TEAM_FILE="teams/backend-team.txt" go run pick_next.go

# Container usage with custom team file
podman run -it --rm \
  -v ./teams:/app/teams \
  -e TEAM_FILE=teams/backend.txt \
  quay.io/asoro/rhdh-install-daily-scrum-picker:latest
```

## Container Registry Setup

### GitLab CI/CD Variables

To enable automatic container builds and pushes to Quay.io, configure these variables in your GitLab project settings (`Settings > CI/CD > Variables`):

| Variable | Description | Example Value |
|----------|-------------|---------------|
| `QUAY_USERNAME` | Your Quay.io username | `your-username` |
| `QUAY_PASSWORD` | Your Quay.io password or robot token | `your-password-or-token` |
| `QUAY_REPOSITORY` | Target repository on Quay.io | `asoro/rhdh-install-daily-scrum-picker` |

### Container Image Tags

The CI pipeline creates different tags based on the git ref:

- **Main branch**: `latest` and `<commit-sha>`
- **Tagged releases**: `<tag-name>` and `latest`
- **Feature branches**: `<branch-name>-<commit-sha>` (manual trigger)

## Technical Details

- **Language**: Go 1.24.4
- **Dependencies**: Standard library + `golang.org/x/term` for raw terminal input
- **Container**: Multi-stage Docker build with Alpine Linux runtime
- **Registry**: Automated builds pushed to Quay.io via GitLab CI
- **Configuration**: Team file (default: `team.txt`, configurable via `TEAM_FILE` env var)
- **Input Mode**: Raw terminal input for immediate keypress response
- **Randomization**: Uses `math/rand` with time-based seeding

## File Structure

```
daily-scrum-picker/
├── .gitlab-ci.yml      # GitLab CI/CD pipeline configuration
├── Dockerfile          # Container image build instructions
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── pick_next.go        # Main application code
├── team.txt            # Team members configuration
└── README.md          # This file
``` 