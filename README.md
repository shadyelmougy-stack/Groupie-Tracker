# Groupie Tracker

## Description

Groupie Tracker is a web application built with **Go (Golang)** that displays information about music artists and bands. It retrieves data from the Groupie Trackers API and presents it through an interactive website.

The application brings together information about artists, their members, their first album, their creation date, concert locations, concert dates, and the relationships between locations and dates.

## Features

- Display artists and bands with their information.
- Show artist details, including:
  - Name
  - Image
  - Members
  - Creation date
  - First album
- Display concert locations for artists.
- Display concert dates.
- Show the relationship between artists, locations, and concert dates.
- Use client-side events and interactions to make browsing easier.
- Handle invalid requests and unexpected errors.

## API

The project uses the Groupie Trackers API:

`https://groupietrackers.herokuapp.com/api`

The API provides the following endpoints:

| Endpoint | Description |
|---|---|
| `/artists` | Artist and band information |
| `/locations` | Concert locations |
| `/dates` | Concert dates |
| `/relation` | Relationships between artists, locations, and dates |

## Technologies Used

- **Go (Golang)** — server-side application and HTTP handling
- **HTML** — page structure
- **CSS** — styling and layout
- **Groupie Trackers API** — artist and concert data

## Getting Started

### Prerequisites

- Go installed on your computer
- Git (if you want to clone the repository)

### Installation

Clone the repository:

```bash
git clone <repository-url>
```

Move into the project directory:

```bash
cd <project-folder>
```

### Run the Application

Start the application:

```bash
go run .
```

Open your browser and visit the local address configured by the application. For example:

`http://localhost:8080`

> If your application uses a different port or start command, update this section accordingly.

## Project Structure

The exact structure may differ depending on your implementation. A typical layout could look like this:

```text
groupie-tracker/
├── main.go
├── go.mod
├── handlers/
├── templates/
├── static/
│   ├── css/
└── README.md
```

## How It Works

1. The Go application starts an HTTP server.
2. The application requests artist, location, date, and relation data from the API.
3. The server processes the retrieved data.
4. HTML templates display the information in the browser.
5. Client-side interactions allow users to explore the available information.

## Error Handling

The application should handle common errors, such as:

- Invalid routes
- Missing or invalid data
- API connection failures
- Internal server errors

## Authors

- Shady Elmougy

## License

This project was created for educational purposes.
