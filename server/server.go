package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"vinyl-player/sonos"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	Spotify "github.com/zmb3/spotify/v2"
	SpotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type ErrorResponse struct {
	code    int
	message string
}

//var config = Configuration{}
var spotify *Spotify.Client
var storage = VinylStorage{
	spotify: spotify,
}
var vaildPaths = regexp.MustCompile("^/(edit|save|view|play)/([0-9]+)$")
var sonosPlayer = sonos.SonosPlayer{}

var ctx = context.Background()

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	sonosPlayer = sonos.New(os.Getenv("SONOS_PLAYER"))
	storage.Connect(os.Getenv("REDIS_URI"))

	spotify = connectSpotify(&clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	})
	storage.spotify = spotify

	server := gin.Default()
	//server.POST("/album", saveAlbum)
	server.GET("/album/:id", getAlbumById)
	server.GET("/album", getAllAlbums)
	server.Static("/assets", "./static/assets")

	// Server always index for SPA
	server.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	server.POST("/barcode", func(c *gin.Context) {
		r, e := CreateFromBarcode(c.PostForm("barcode"))
		if e != nil {
			c.AbortWithStatusJSON(404, gin.H{"message": e.Error()})
			return
		}

		if c.PostForm("save") == "1" {
			storage.Create(r)
		}

		c.IndentedJSON(200, r)
	})

	server.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")))
}

func connectSpotify(cfg *clientcredentials.Config) *Spotify.Client {
	cfg.TokenURL = SpotifyAuth.TokenURL

	token, err := cfg.Token(ctx)
	if err != nil {
		panic(err)
	}

	httpClient := SpotifyAuth.New().Client(ctx, token)
	return Spotify.New(httpClient)
}

func saveAlbum(c *gin.Context) {
	newAlbum := VinylAlbum{}

	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	if _, err := storage.Create(&newAlbum); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	a, e := storage.getOne(id, false)
	if e != nil {
		c.AbortWithError(http.StatusNotFound, e)
		return
	}

	c.IndentedJSON(http.StatusOK, a)
}

func getAllAlbums(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, storage.getAll())
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	id := vaildPaths.FindStringSubmatch(r.URL.Path)

	if id == nil {
		http.NotFound(w, r)
		return
	}

	//idInt, _ := strconv.ParseInt(id[2], 10, 8)
	v, e := storage.getOne(id[2], true)
	if e != nil {
		http.Error(w, e.Error(), 404)
		return
	}

	e = sonosPlayer.PlaySpotifyAlbum(v.Links.SpotifyAlbumURI)
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	http.Redirect(w, r, "/overview", 307)
}
