# Rate My

A simple web service that allows users to submit 1-5 star ratings for various events or occurrences.

## Features

- Single page web interface with star rating selection (1-5 stars)
- REST API endpoint for submitting ratings
- Logging to stdout and file (`ratings.log`)

## Getting Started

### Prerequisites

- Go 1.22 or later

### Installation

```bash
git clone https://github.com/its-the-vibe/rate-my.git
cd rate-my
```

### Running the Server

```bash
go run main.go
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

### Usage

Open your browser and navigate to `http://localhost:8080` to access the rating interface.

## API

### POST /rate

Submit a rating for an event.

**Request Body:**

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "event": "tube journey",
  "rating": 4
}
```

| Field | Type | Description |
|-------|------|-------------|
| timestamp | string | ISO 8601 timestamp (optional, defaults to current time) |
| event | string | The event being rated (e.g., "tube journey") |
| rating | integer | Rating from 1 to 5 |

**Response:**

```json
{
  "status": "success"
}
```

## Logging

Ratings are logged in two ways:
1. Standard output (stdout)
2. Appended to `ratings.log` file

Log format:
```
[2024-01-15T10:30:00Z] Event: tube journey, Rating: 4
```

## Project Structure

```
rate-my/
├── main.go          # Go web server
├── static/
│   └── index.html   # Single page web interface
├── go.mod           # Go module file
├── ratings.log      # Rating log file (created on first submission)
└── README.md        # This file
```

## Future Enhancements

- Dropdown to select different events to rate
- Date picker for rating past events
- Publish ratings to pub/sub
- Persistent storage

## License

MIT
