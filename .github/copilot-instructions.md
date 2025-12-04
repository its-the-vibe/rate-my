# Copilot Instructions for Rate My

## Project Overview

Rate My is a simple Go web service that allows users to submit 1-5 star ratings for events. It features a single-page web interface with star rating selection and a REST API endpoint for submitting ratings.

## Technology Stack

- **Language**: Go 1.22+
- **Web Framework**: Standard library (`net/http`)
- **Frontend**: Vanilla HTML, CSS, and JavaScript (no frameworks)
- **Storage**: File-based logging (`ratings.log`)

## Commands

### Build

```bash
go build -o rate-my .
```

### Run

```bash
go run main.go
```

Or with a custom port:

```bash
PORT=3000 go run main.go
```

### Test

```bash
go test ./...
```

### Format

```bash
gofmt -w .
```

### Lint

```bash
go vet ./...
```

## Project Structure

```
rate-my/
├── main.go          # Go web server with HTTP handlers
├── static/
│   └── index.html   # Single page web interface with star rating UI
├── go.mod           # Go module definition
├── ratings.log      # Rating log file (created at runtime, gitignored)
└── README.md        # Project documentation
```

## Code Style

### Go

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting
- Use meaningful variable and function names
- Keep functions focused and concise
- Handle errors explicitly - do not ignore returned errors
- Use standard library packages when possible

### HTML/CSS/JavaScript

- Use semantic HTML elements
- Keep CSS in `<style>` blocks within the HTML file (single-file approach)
- Use vanilla JavaScript - no external frameworks
- Use `const` and `let` instead of `var`
- Use async/await for asynchronous operations

## API Reference

### POST /rate

Submit a rating for an event.

**Request Body:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",  // optional, defaults to current time
  "event": "tube journey",
  "rating": 4  // must be 1-5
}
```

**Response:**
```json
{
  "status": "success"
}
```

## Git Workflow

- Write clear, descriptive commit messages
- Keep commits focused on a single change
- Test changes before committing

## Boundaries and Constraints

### Do Not

- Add external dependencies without careful consideration
- Modify the `.gitignore` patterns unless necessary
- Commit `ratings.log` or other log files
- Introduce breaking changes to the API without updating documentation

### Security

- Never commit secrets or credentials
- Validate all user input
- Use appropriate HTTP status codes for errors
