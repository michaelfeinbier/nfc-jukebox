package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v9"
)

type uuid = string
type VinylAlbum struct {
	Id     int64 `redis:"Id"`
	Name   string
	Artist string
	Links  AlbumLinks
}

type AlbumLinks struct {
	SpotifyURI    string `json:"spotify"`
	MusicBrainzId uuid   `json:"mbid"`
}

type VinylStorage struct {
	Albums []VinylAlbum
	redis  *redis.Client
}

var ctx = context.Background()

const (
	ALBUM_LIST_KEY   = "album_list"
	ALBUM_KEY_PREFIX = "album"
	LINK_KEY_PREFIX  = "link"
)

//var client

func (s *VinylStorage) Connect(uri string) {
	opt, _ := redis.ParseURL(uri)
	s.redis = redis.NewClient(opt)
}

func (s *VinylStorage) getAll() []VinylAlbum {
	r := s.redis.SMembers(ctx, ALBUM_LIST_KEY)

	all := s.redis.MGet(ctx, r.Val()...)
	var data = []VinylAlbum{}

	for _, r := range all.Val() {
		a := VinylAlbum{}
		if e := json.Unmarshal([]byte(r.(string)), &a); e != nil {
			panic(e)
		}
		data = append(data, a)
	}

	return data
}

func (s *VinylStorage) getOne(ID string) (VinylAlbum, error) {
	key := fmt.Sprintf("%s:%s", ALBUM_KEY_PREFIX, ID)
	res := s.redis.Get(ctx, key)

	if res.Err() != nil {
		return VinylAlbum{}, res.Err()
	}

	var album = VinylAlbum{}
	if err := json.Unmarshal([]byte(res.Val()), &album); err != nil {
		return album, err
	}

	return album, nil
}

func (s *VinylStorage) Create(album *VinylAlbum) (*VinylAlbum, error) {
	album.Id = s.getNewId()
	r, k, e := s.Save(album)
	if e != nil {
		return album, e
	}

	s.redis.SAdd(ctx, ALBUM_LIST_KEY, k)

	return r, nil
}

func (s *VinylStorage) Save(album *VinylAlbum) (cAlbum *VinylAlbum, key string, e error) {
	if album.Id == 0 {
		e = fmt.Errorf("Record does not have ID")
		return
	}

	ja, _ := json.Marshal(album)
	key = fmt.Sprintf("%s:%d", ALBUM_KEY_PREFIX, album.Id)
	r := s.redis.Set(ctx, key, ja, 0)
	if r.Err() != nil {
		e = r.Err()
		return
	}

	cAlbum = album

	return
}

func (s *VinylStorage) getNewId() int64 {
	return s.redis.Incr(ctx, "next_album_id").Val()
}
