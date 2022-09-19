package main

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
)

func getMusicBrainzRelease(barcode string) ([]*gomusicbrainz.Release, error) {
	mbClient, e := gomusicbrainz.NewWS2Client("https://musicbrainz.org/ws/2", "Michael's Vinyl Database", "0.0.1", "https://github.com/michaelfeinbier/go-vinyl-nfc-playback")
	if e != nil {
		return []*gomusicbrainz.Release{}, e
	}

	resp, e := mbClient.SearchRelease(fmt.Sprintf("barcode:%s", barcode), 1, 0)
	if e != nil {
		return []*gomusicbrainz.Release{}, e
	}

	if resp.Count == 0 {
		return []*gomusicbrainz.Release{}, fmt.Errorf(`Did not find barcode "%s" in musicbrainz`, barcode)
	}

	return resp.Releases, nil
}
