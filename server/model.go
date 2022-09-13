package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type VinylAlbum struct {
	Id     int64
	Name   string
	Artist string
	Links  AlbumLinks
}

type AlbumLinks struct {
	SpotifyURI string `json:"spotify"`
}

type VinylStorage struct {
	Albums   []VinylAlbum
	fileName string
}

func (s *VinylStorage) getAll() []VinylAlbum {
	return s.readJSON(s.fileName, func(va VinylAlbum) bool { return true }, false)
}

func (s *VinylStorage) getOne(ID int64) (VinylAlbum, error) {
	res := s.readJSON(s.fileName, func(va VinylAlbum) bool {
		return va.Id == ID
	}, true)

	if len(res) == 0 {
		return VinylAlbum{}, fmt.Errorf("Did not find album with id %d in storage", ID)
	}

	return res[0], nil
}

func (s *VinylStorage) readJSON(fileName string, filter func(VinylAlbum) bool, onlyOne bool) []VinylAlbum {
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := json.NewDecoder(file)

	filteredData := []VinylAlbum{}

	decoder.Token()

	data := VinylAlbum{}
	for decoder.More() {
		decoder.Decode(&data)

		if filter(data) {
			filteredData = append(filteredData, data)

			if onlyOne {
				return filteredData
			}
		}
	}

	return filteredData
}
