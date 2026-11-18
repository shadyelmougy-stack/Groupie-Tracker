package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

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
