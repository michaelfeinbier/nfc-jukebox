package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Discogs struct {
	Token string
}

type DiscogsSearchResult struct {
	Pagination *DiscogsPagination `json:"pagination"`
	Results    []*DiscogsAlbum    `json:"results"`
}

type DiscogsPagination struct {
	Page    int `json:"page"`
	Pages   int `json:"pages"`
	PerPage int `json:"per_page"`
	Items   int `json:"items"`
}

type DiscogsAlbum struct {
	Id            int      `json:"id"`
	Title         string   `json:"title"`
	Type          string   `json:"type"`
	Format        []string `json:"format"`
	Year          string   `json:"year"`
	MasterId      int      `json:"master_id"`
	CatalogNumber string   `json:"catno"`
	Image         string   `json:"cover_image"`
	Thumbnail     string   `json:"thumb"`
}

const (
	DISCOGS_BASE = "https://api.discogs.com"
)

func (d *Discogs) FindByQuery(q string) (DiscogsSearchResult, error) {
	res := DiscogsSearchResult{}

	req, e := http.NewRequest("GET", fmt.Sprintf("%s/database/search?q=%s&type=release", DISCOGS_BASE, q), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Discogs token=%s", d.Token))
	if e != nil {
		return res, e
	}

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return res, e
	}

	defer resp.Body.Close()

	if e := json.NewDecoder(resp.Body).Decode(&res); e != nil {
		return res, e
	}

	return res, nil
}
