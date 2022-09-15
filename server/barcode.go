package main

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
	Spotify "github.com/zmb3/spotify/v2"
)

// Create / Enrich a Record based on metadata fetched via barcode

func CreateFromBarcode(barcode string) (a *VinylAlbum, e error) {
	mbr, e := getMusicBrainzRelease(barcode)
	if e != nil {
		return &VinylAlbum{}, e
	}

	a = &VinylAlbum{
		Name: mbr.Title,
		Links: AlbumLinks{
			MusicBrainzId: string(mbr.ID),
		},
		Artist: mbr.ArtistCredit.NameCredits[0].Artist.Name,
		Metadata: AlbumMetadata{
			UPCRelease: mbr.Barcode,
		},
	}

	s := lookupSpotifyUri(a.Name, a.Artist)
	a.Links.SpotifyAlbumURI = string(s.URI)
	a.Links.SpotifyArtistURI = string(s.Artists[0].URI)

	storage.GetSpotifyMetadata(s.ID, &a.Metadata)

	return
}

func getMusicBrainzRelease(barcode string) (*gomusicbrainz.Release, error) {
	mbClient, e := gomusicbrainz.NewWS2Client("https://musicbrainz.org/ws/2", "Michael's Vinyl Database", "0.0.1", "https://github.com/michaelfeinbier/go-vinyl-nfc-playback")
	if e != nil {
		return &gomusicbrainz.Release{}, e
	}

	resp, e := mbClient.SearchRelease(fmt.Sprintf("barcode:%s", barcode), 1, 0)
	if e != nil {
		return &gomusicbrainz.Release{}, e
	}

	if resp.Count == 0 {
		return &gomusicbrainz.Release{}, fmt.Errorf(`Did not find barcode "%s" in musicbrainz`, barcode)
	}

	return resp.Releases[0], nil
}

func lookupSpotifyUri(title string, artist string) Spotify.SimpleAlbum {
	q := fmt.Sprintf("artist:%s album:%s", artist, title)
	r, _ := spotify.Search(ctx, q, Spotify.SearchTypeAlbum)

	if r.Albums.Total > 0 {
		return r.Albums.Albums[0]
	}

	return Spotify.SimpleAlbum{}
}
