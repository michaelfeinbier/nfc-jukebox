package main

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
	Spotify "github.com/zmb3/spotify/v2"
)

type GoldenRecord struct {
}

type CombinedSearchResult struct {
	Name        string                   `json:"name"`
	Artist      string                   `json:"artist"`
	MusicBrainz []*gomusicbrainz.Release `json:"musicbrainz"`
	Discogs     []*DiscogsAlbum          `json:"discogs"`
	Spotify     *[]Spotify.SimpleAlbum   `json:"spotify"`
}

func (g *GoldenRecord) CombinedSearch(q string) (CombinedSearchResult, error) {
	cr := CombinedSearchResult{}

	// fetch discogs first
	dr, _ := discogs.FindByQuery(q)

	cr.Discogs = dr.Results

	// musicbrainz
	mr, _ := getMusicBrainzRelease(q)
	cr.MusicBrainz = mr

	g.GuessSpotifyFromResults(&cr)
	g.GuessTitleFromResults(&cr)

	return cr, nil
}

func (g *GoldenRecord) GuessSpotifyFromResults(s *CombinedSearchResult) {

	var query string
	if len(s.MusicBrainz) > 0 {
		r := s.MusicBrainz[0]
		query = fmt.Sprintf("artist:%s album:%s", r.ArtistCredit.NameCredits[0].Artist.Name, r.Title)
	} else if len(s.Discogs) > 0 {
		r := s.Discogs[0]
		query = r.Title
	} else {
		return
	}

	r, _ := spotify.Search(ctx, query, Spotify.SearchTypeAlbum)

	if r.Albums.Total > 0 {
		s.Spotify = &r.Albums.Albums
	}
}

func (g *GoldenRecord) GuessTitleFromResults(s *CombinedSearchResult) error {
	if len(s.MusicBrainz) > 0 {
		r := s.MusicBrainz[0]
		s.Artist = r.ArtistCredit.NameCredits[0].Artist.Name
		s.Name = r.Title
		return nil
	} else if len(*s.Spotify) > 0 {
		r := (*s.Spotify)[0]
		s.Artist = r.Artists[0].Name
		s.Name = r.Name
		return nil
		//} else if len(s.Discogs) > 0 {
		//	return nil
	} else {
		return fmt.Errorf("Could not create Title from search results")
	}
}
