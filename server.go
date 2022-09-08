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
}

func main() {
	config := Configuration{}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal().Err(err)
	}

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal().Err(err)
	}

	//http.HandleFunc("/mytest", indexHandler)

	log.Info().Msgf("Starting Server on Port %d", config.Port)
	server :=
		http.ListenAndServe(fmt.Sprintf(":%d", config.Port),
			http.FileServer(http.Dir("static/")))
	log.Fatal().Err(server)

}
