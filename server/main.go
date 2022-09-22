package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"vinyl-player/sonos"

	"github.com/gin-gonic/gin"
	Spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
	SpotifyAuth "golang.org/x/oauth2/spotify"
)

type ErrorResponse struct {
	code    int
	message string
}

var spotify *Spotify.Client
var storage = VinylStorage{
	spotify: spotify,
}
var sonosPlayer = sonos.SonosPlayer{}

var ctx = context.Background()
var discogs *Discogs
var goldenRecord *GoldenRecord

func main() {
	sonosPlayer = sonos.New(&sonos.SonosConfig{
		IpAddress:     os.Getenv("SONOS_PLAYER"),
		CoordinatorId: os.Getenv("SONOS_COORDINATOR"),
	})
	storage.Connect(os.Getenv("REDIS_URI"))
	discogs = &Discogs{
		Token: os.Getenv("DISCOGS_TOKEN"),
	}

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
	server.GET("/p/:id", playAlbum)

	// API routes
	api := server.Group("/api")
	{
		api.POST("/album", createAlbum)
		api.GET("/album/:id", getAlbumById)
		api.GET("/album", func(c *gin.Context) {
			c.IndentedJSON(http.StatusOK, storage.getAll())
		})

		// connect services to search
		search := api.Group("/search")
		{
			search.GET("/:q", searchAll)
			search.GET("/discogs/:q", searchDiscogs)
		}
	}

	server.Run(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")))
}
