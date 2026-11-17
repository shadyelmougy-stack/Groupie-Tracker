package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const apiURL = "http://localhost:8081/api"

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type RelationIndex struct {
	Index []Relation `json:"index"`
}

type ArtistPage struct {
	Artist    Artist
	Locations map[string][]string
}

type HomePage struct {
	Query   string
	Artists []Artist
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	var artists []Artist
	err := fetchJSON(apiURL+"/artists", &artists)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not load artists", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query().Get("q")
	page := HomePage{
		Query:   query,
		Artists: filterArtists(artists, query),
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid artist id", http.StatusBadRequest)
		return
	}

	var artists []Artist
	err = fetchJSON(apiURL+"/artists", &artists)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not load artists", http.StatusInternalServerError)
		return
	}

	artist, found := findArtist(artists, id)
	if !found {
		http.NotFound(w, r)
		return
	}

	var relationIndex RelationIndex
	err = fetchJSON(apiURL+"/relation", &relationIndex)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not load tour data", http.StatusInternalServerError)
		return
	}

	relation, _ := findRelation(relationIndex.Index, artist.ID)

	page := ArtistPage{
		Artist:    artist,
		Locations: relation.DatesLocations,
	}

	funcs := template.FuncMap{"formatPlace": formatPlace}
	tmpl, err := template.New("artist.html").Funcs(funcs).ParseFiles("templates/artist.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}

func findRelation(relations []Relation, id int) (Relation, bool) {
	for _, r := range relations {
		if r.ID == id {
			return r, true
		}
	}
	return Relation{}, false
}

func findArtist(artists []Artist, id int) (Artist, bool) {
	for _, a := range artists {
		if a.ID == id {
			return a, true
		}
	}
	return Artist{}, false
}

var client = &http.Client{Timeout: 10 * time.Second}

func fetchJSON(url string, target any) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fetching %s: unexpected status: %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func formatPlace(raw string) string {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return raw
	}

	city := titleCase(strings.ReplaceAll(parts[0], "_", " "))

	country := strings.ReplaceAll(parts[1], "_", " ")
	if len(country) <= 3 {
		country = strings.ToUpper(country)
	} else {
		country = titleCase(country)
	}

	return city + ", " + country
}

func titleCase(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}

func filterArtists(artists []Artist, query string) []Artist {
	query = strings.ToLower(strings.TrimSpace(query))

	var result []Artist
	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), query) {
			result = append(result, a)
		}
	}
	return result
}
