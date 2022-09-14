package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"vinyl-player/sonos"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	Port         int
	ClientKey    string
	ClientSecret string
	SonosBaseURI string
	RedisURI     string
}

var config = Configuration{}
var storage = VinylStorage{}
var vaildPaths = regexp.MustCompile("^/(edit|save|view|play)/([0-9]+)$")
var sonosPlayer sonos.SonosPlayer

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal().Err(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal().Err(err)
	}

	sonosPlayer = sonos.New("192.168.1.44")
	storage.Connect(config.RedisURI)

	server := gin.Default()
	server.POST("/album", saveAlbum)
	server.GET("/album/:id", getAlbumById)
	server.GET("/album", getAllAlbums)
	server.Static("/assets", "./static/assets")
	server.StaticFile("/", "./static/index.html")

	server.Run(fmt.Sprintf(":%d", config.Port))
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

	a, e := storage.getOne(id)
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
	v, e := storage.getOne(id[2])
	if e != nil {
		http.Error(w, e.Error(), 404)
		return
	}

	e = sonosPlayer.PlaySpotifyAlbum(v.Links.SpotifyURI)
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	http.Redirect(w, r, "/overview", 307)
}
