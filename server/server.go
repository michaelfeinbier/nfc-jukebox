package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"vinyl-player/sonos"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	Port         int
	ClientKey    string
	ClientSecret string
	SonosBaseURI string
}

var config = Configuration{}
var storage = VinylStorage{
	fileName: "album.json",
}
var templates = template.Must(template.ParseGlob("templates/*.html"))
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

	//http.HandleFunc("/auth", startAuthentication)
	//http.HandleFunc("/auth/redirect", handleCallback)
	http.HandleFunc("/overview", handleIndex)
	http.HandleFunc("/play/", handlePlay)
	http.Handle("/", http.FileServer(http.Dir("static/")))

	log.Info().Msgf("Starting Server on Port %d", config.Port)
	server := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal().Err(server)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", storage.getAll())
}

func handlePlay(w http.ResponseWriter, r *http.Request) {
	id := vaildPaths.FindStringSubmatch(r.URL.Path)

	if id == nil {
		http.NotFound(w, r)
		return
	}

	idInt, _ := strconv.ParseInt(id[2], 10, 8)
	v, e := storage.getOne(idInt)
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

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
