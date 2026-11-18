# Groupie Tracker

A small website written in Go that shows information about bands and artists: members, start year, first album, and concert places and dates.

The data comes from a REST API (the Groupie Trackers API). The website is built using **only the Go standard library**.

## Features

- **Home page**: all artists shown as cards (picture, name, start year, first album, members).
- **Search box (the client-server event)**: the visitor types part of an artist's name and presses *Search*. The browser sends `GET /?q=...` to the server, and the server answers with only the matching artists. The search is case-insensitive, and an empty search shows everyone.
- **Details page** for each artist: members and a list of concerts (places and dates). Place names are shown in a readable form, for example `north_carolina-usa` becomes `North Carolina, USA`.
- **Error handling**: wrong requests, unknown pages and API failures receive a proper HTTP error response. The server does not crash.
- **Unit tests** for the helper functions.

## Requirements

- Go 1.18 or newer
- The Groupie Trackers API running locally at `http://localhost:8081/api`

## How to run

1. Start the API (as given in the project instructions) so that `http://localhost:8081/api` answers.
2. In the project folder, run:

   ```
   go run .
   ```

3. Open <http://localhost:8080> in your browser.

> Use `go run .` and not `go run main.go`: the program is made of several files, and `go run main.go` would compile only one of them.

## Project structure

```
groupie-tracker/
├── go.mod
├── main.go            # starts the server, routes, and the HTTP handlers
├── data.go            # data types, API address, and the code that fetches and searches API data
├── helpers.go         # small helpers: formatPlace, titleCase, filterArtists
├── helpers_test.go    # unit tests for the helpers
├── static/
│   └── style.css      # styling (served under /static/)
└── templates/
    ├── index.html     # home page template
    └── artist.html    # artist details template
```

## Routes

| Route | Description |
|---|---|
| `GET /` | Home page with all artists. Optional query `?q=text` filters artists by name. |
| `GET /artist?id=N` | Details page of the artist whose `id` is `N`. |
| `GET /static/...` | Static files (the stylesheet). |

## Error handling

| Situation | Response |
|---|---|
| Missing or invalid artist id (`/artist?id=abc`, `/artist`) | `400 Bad Request` |
| Unknown path, or an artist id that does not exist | `404 Not Found` |
| The API is unreachable, too slow (10 second timeout), answers with an error status, or sends invalid JSON | `500 Internal Server Error` with a short message; the real error is written to the server log |
| A template file cannot be loaded | `500 Internal Server Error` |
| An artist exists but has no tour data | The page is shown with an empty concerts list |

## Tests

```
go test ./...
```

The tests cover `formatPlace` (readable place names) and `filterArtists` (the search logic).

## Configuration

- API address: the `apiURL` constant in `data.go`.
- API request timeout: the `client` variable in `data.go` (10 seconds).
- Server port: `:8080` in `main.go`.

## About the API data

- `/artists` returns a list of artists (id, image, name, members, creation date, first album).
- `/relation` returns all tour data at once, wrapped as `{"index": [{"id": 1, "datesLocations": {...}}, ...]}`. Each entry maps a place to a list of dates. The site finds an artist's entry by matching its `id`.
