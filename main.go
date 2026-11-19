package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	indexTmpl  *template.Template
	artistTmpl *template.Template
)

func main() {
	err := loadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running at http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func loadTemplates() error {
	var err error

	indexTmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		return err
	}

	funcs := template.FuncMap{"formatPlace": formatPlace}
	artistTmpl, err = template.New("artist.html").Funcs(funcs).ParseFiles("templates/artist.html")
	return err
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	err = indexTmpl.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	err = artistTmpl.Execute(w, page)
	if err != nil {
		log.Println(err)
	}
}
