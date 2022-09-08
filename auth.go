package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

type SonosAuthentication struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

const (
	SONOS_LOGIN = "https://api.sonos.com/login/v3"
	SCOPE       = "playback-control-all"
	STATE       = "demo"
)

func startAuthentication(w http.ResponseWriter, r *http.Request) {
	p := url.Values{}
	p.Add("client_id", config.ClientKey)
	p.Add("redirect_uri", getRedirectURL())
	p.Add("response_type", "code")
	p.Add("scope", SCOPE)
	p.Add("state", STATE)
	redirectUrl := fmt.Sprintf("%s/oauth?%s", SONOS_LOGIN, p.Encode())

	log.Debug().Msg("Redirect to: " + redirectUrl)

	http.Redirect(w, r, redirectUrl, 307)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v.Get("state") == STATE && len(v.Get("code")) > 0 {
		log.Debug().
			Str("exchangeToken", v.Get("code")).
			Msg("Got valid exchange token")

		form := url.Values{}
		form.Add("grant_type", "authorization_code")
		form.Add("code", v.Get("code"))
		form.Add("redirect_uri", getRedirectURL())

		log.Info().
			Str("url", form.Encode()).
			Msg("encode form")

		// exchange request
		ex, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/access", SONOS_LOGIN), strings.NewReader(form.Encode()))
		ex.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		ex.SetBasicAuth(config.ClientKey, config.ClientSecret)

		if err != nil {
			log.Error().Err(err)
			return
		}

		resp, err := http.DefaultClient.Do(ex)
		if err != nil {
			log.Error().Err(err)
			http.Error(w, err.Error(), 500)
			return
		}

		//defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			log.Warn().Int("HTTP Status", resp.StatusCode)
			fmt.Fprintf(w, "Something went wrong (%d): %s", resp.StatusCode, string(body))
			return
		}

		writeToken(body)
		fmt.Fprintf(w, "OK"+string(body))
		return
	} else {
		http.Error(w, "Invalid Params", http.StatusBadRequest)
		return
	}
}

func validAccessToken() bool {
	return false
}

func getRedirectURL() string {
	return fmt.Sprintf("http://localhost:%d/auth/redirect", config.Port)
}

func writeToken(body []byte) (token SonosAuthentication, e error) {
	token = SonosAuthentication{}
	e = json.Unmarshal(body, &token)

	if e != nil {
		fmt.Println("Error:", e)
		return
	}
	e = ioutil.WriteFile("token.json", body, 0644)
	return
}
