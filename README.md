# Daily Scrum Picker

[![Tests](https://github.com/rm3l/daily-scrum-picker/actions/workflows/test.yaml/badge.svg)](https://github.com/rm3l/daily-scrum-picker/actions/workflows/test.yaml)
<!--[![codecov](https://codecov.io/gh/rm3l/daily-scrum-picker/branch/main/graph/badge.svg)](https://codecov.io/gh/rm3l/daily-scrum-picker)-->
[![Build and publish container image](https://github.com/rm3l/daily-scrum-picker/actions/workflows/build-and-publish-container-image.yaml/badge.svg)](https://github.com/rm3l/daily-scrum-picker/actions/workflows/build-and-publish-container-image.yaml)

A simple Go utility to fairly select the next person to speak during daily scrum/stand-up meetings.

**NOTE**: Assisted by AI as a quick experiment to get something up and running in a few minutes.

## Overview

This tool ensures fair rotation of team members during daily stand-ups by automatically tracking who has spoken and resetting when everyone has had a turn.

**Features:**

1. **Fair Rotation**: Tracks who hasn't spoken yet automatically
2. **Random Shuffling**: When everyone has had a turn, shuffles the team list for the next cycle
3. **Persistent Tracking**: Remembers selections between runs so you can use it daily
4. **Automatic Reset**: When the list is empty, automatically starts a new randomized cycle

## Installation

### Prerequisites

- Go 1.23+ (tested with Go 1.23 and 1.24)

### Setup

```bash
git clone https://github.com/rm3l/daily-scrum-picker.git && cd daily-scrum-picker
```

## Usage

### Local Development

1. Create a local `team.txt` file containing the names of your team members. See [Team members](#team-members) for more details.

2. Simply run with Go:

```bash
go run pick_next.go
```

Alternatively, you can build and run the executable:

```bash
go build -o daily-scrum-picker pick_next.go
./daily-scrum-picker
```

You can also specify a custom team file using the `--team-file` (or `-t`) flag:

```bash
go run pick_next.go --team-file=/path/to/my-team.txt
go run pick_next.go -t /path/to/my-team.txt
./daily-scrum-picker --team-file=teams/backend.txt
./daily-scrum-picker -t teams/backend.txt
```

Or read team members from stdin using `-`:

```bash
echo -e "Alice\nBob\nCharlie" | ./daily-scrum-picker -t -
cat team-members.txt | go run pick_next.go --team-file=-
```

### Container Usage

This tool is also available as a container image on ghcr.io. This allows you to use the tool without cloning the repository or installing Go.

It runs in interactive mode with **single-keypress commands** - no need to press Enter:

```bash
# Mount your custom team file
podman run -it --rm \
  -v ./my-team.txt:/app/team.txt \
  ghcr.io/rm3l/daily-scrum-picker:main

# Or use a different file path with environment variable
podman run -it --rm \
  -v ./teams:/app/teams \
  -e TEAM_FILE=teams/backend.txt \
  ghcr.io/rm3l/daily-scrum-picker:main
```

### Interactive Mode

**Available commands (single keypress):**

- **`p`** - Pick the next person for daily scrum
- **`r`** - Reset and start over with all team members  
- **`s`** - Show current status and remaining team members
- **`h`** - Show help message
- **`q`** - Exit the program

**Notes:** 

- Use the `-it` flags to enable interactive mode with proper terminal support
- Commands respond immediately without pressing Enter
- Fallback to Enter-required mode if raw terminal access is unavailable

### Output Examples

**Interactive session:**

```txt
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

Team members are configured via a team file or stdin. By default, the tool looks for `team.txt` in the same directory where it is run, but you can specify a different source using either:

1. The `--team-file` (or `-t`) command-line flag (file path or `-` for stdin)
2. The `TEAM_FILE` environment variable

Create or edit the team file with one team member name per line:

```txt
# Daily Scrum Team Members
# Add one team member name per line
# Lines starting with # are ignored

Alice
Bob
Charlie
Diana
```

#### Configuration Options

The tool supports multiple ways to specify the team source (in order of precedence):

| Method | Description | Priority |
|--------|-------------|----------|
| `--team-file` / `-t` flag | Path to team members file or `-` for stdin | **Highest** (overrides environment variable) |
| `TEAM_FILE` environment variable | Path to team members file | Medium |
| Default | Uses `team.txt` in current directory | Lowest |

**Examples:**

```bash
# Using command-line flag (long form)
go run pick_next.go --team-file="/path/to/my-team.txt"
./daily-scrum-picker --team-file="teams/backend.txt"

# Using command-line flag (short form)
go run pick_next.go -t "/path/to/my-team.txt"
./daily-scrum-picker -t "teams/backend.txt"

# Reading from stdin
echo -e "Alice\nBob\nCharlie\nDiana" | ./daily-scrum-picker -t -
cat my-team.txt | go run pick_next.go --team-file=-

# Using environment variable
export TEAM_FILE="/path/to/my-team.txt"
go run pick_next.go

# Environment variable for single run
TEAM_FILE="/path/to/teams/backend-team.txt" go run pick_next.go

# Command-line flag takes precedence over environment variable
TEAM_FILE="/path/to/env-team.txt" go run pick_next.go -t "/path/to/flag-team.txt"
# Will use /path/to/flag-team.txt (flag overrides environment variable)

# Stdin takes precedence over environment variable
TEAM_FILE="/path/to/env-team.txt" echo -e "User1\nUser2" | ./daily-scrum-picker -t -
# Will read from stdin (flag overrides environment variable)

# Container usage with environment variable
podman run -it --rm \
  -v ./teams:/app/teams \
  -e TEAM_FILE=teams/backend.txt \
  ghcr.io/rm3l/daily-scrum-picker:main

# Container usage with custom team file mounted
podman run -it --rm \
  -v ./my-team.txt:/app/team.txt \
  ghcr.io/rm3l/daily-scrum-picker:main -t /app/team.txt

# Dynamic team selection from command output
git log --format="%an" --since="1 month ago" | sort -u | ./daily-scrum-picker -t -
# Uses recent Git contributors as team members

# Integration with other tools
curl -s https://api.github.com/repos/owner/repo/contributors | jq -r '.[].login' | head -10 | ./daily-scrum-picker -t -
# Uses GitHub contributors as team members
```

## Development

### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -race -coverprofile=coverage.out
go tool cover -func=coverage.out

# Run benchmarks
go test -bench=. -benchmem

# Run linting (requires golangci-lint)
golangci-lint run
```

### CI/CD

The project includes comprehensive GitHub Actions workflows:

- **Tests** (`test.yaml`): Runs on every push and PR
  - Tests on Go 1.23 and 1.24
  - Includes race detection, benchmarks, and coverage reporting
  - Uploads coverage reports to Codecov
  - Runs linting with golangci-lint
  
- **Build** (`build-and-publish-container-image.yaml`): Builds and publishes container images
  - Runs tests before building to ensure quality
  - Publishes to GitHub Container Registry
  - Signs images with cosign

## License

    The MIT License (MIT)

    Copyright (c) 2025 Armel Soro

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
