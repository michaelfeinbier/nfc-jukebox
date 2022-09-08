package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

	http.HandleFunc("/auth", startAuthentication)
	http.HandleFunc("/auth/redirect", handleCallback)
	http.Handle("/", http.FileServer(http.Dir("static/")))

	log.Info().Msgf("Starting Server on Port %d", config.Port)
	server := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	log.Fatal().Err(server)

}
