package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"vinyl-player/sonos"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	Spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
	SpotifyAuth "golang.org/x/oauth2/spotify"
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
var sonosPlayer = sonos.SonosPlayer{}

var ctx = context.Background()

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	sonosPlayer = sonos.New(&sonos.SonosConfig{
		IpAddress:     os.Getenv("SONOS_PLAYER"),
		CoordinatorId: os.Getenv("SONOS_COORDINATOR"),
	})
	storage.Connect(os.Getenv("REDIS_URI"))

	scfg := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     SpotifyAuth.Endpoint.TokenURL,
	}
	spotify = Spotify.New(scfg.Client(ctx))
	storage.spotify = spotify

	server := gin.Default()

	// Serve index and assets for SPA
	server.Static("/assets", "./static/assets")
	server.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// This is the route coded to the NFC tags
	server.GET("/p/:id", func(c *gin.Context) {
		id := c.Param("id")
		v, _ := storage.getOne(id, false)
		if len(v.Links.SpotifyAlbumURI) > 0 {
			if e := sonosPlayer.PlaySpotifyAlbum(v.Links.SpotifyAlbumURI); e != nil {
				c.AbortWithError(500, e)
			}
		}
		// redirect to SPA app
		c.Redirect(http.StatusTemporaryRedirect, "/view/"+id)
	})

	// API routes
	api := server.Group("/api")
	{
		api.GET("/album/:id", getAlbumById)
		api.GET("/album", func(c *gin.Context) {
			c.IndentedJSON(http.StatusOK, storage.getAll())
		})
		api.POST("/barcode", func(c *gin.Context) {
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
	}

	server.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")))
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
