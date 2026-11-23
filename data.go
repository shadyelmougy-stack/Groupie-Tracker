package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

const apiURL = "http://localhost:8081/api"

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

	return json.NewDecoder(resp.Body).Decode(target)
}

func findArtist(artists []Artist, id int) (Artist, bool) {
	for _, a := range artists {
		if a.ID == id {
			return a, true
		}
	}
	return Artist{}, false
}

func findRelation(relations []Relation, id int) (Relation, bool) {
	for _, r := range relations {
		if r.ID == id {
			return r, true
		}
	}
	return Relation{}, false
}
