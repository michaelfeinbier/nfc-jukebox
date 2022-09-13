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
	SONOS_LOGIN   = "https://api.sonos.com/login/v3"
	SPOTIFY_LOGIN = "https://accounts.spotify.com"
	SCOPE         = "user-modify-playback-state user-read-playback-state user-read-currently-playing"
	STATE         = "demo"
)

func startAuthentication(w http.ResponseWriter, r *http.Request) {
	p := url.Values{}
	p.Add("client_id", config.ClientKey)
	p.Add("redirect_uri", getRedirectURL())
	p.Add("response_type", "code")
	p.Add("scope", SCOPE)
	p.Add("state", STATE)
	redirectUrl := fmt.Sprintf("%s/authorize?%s", SPOTIFY_LOGIN, p.Encode())

	log.Debug().
		Str("redir_url", getRedirectURL()).
		Msg("Redirect to: " + redirectUrl)

	http.Redirect(w, r, redirectUrl, 307)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	if v.Get("state") == STATE && len(v.Get("code")) > 0 {
		log.Debug().
			Str("exchangeToken", v.Get("code")).
			Msg("Got valid exchange token")

		token, err := exchangeAuthCode(v.Get("code"))

		if err != nil {
			http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		saveToken(token)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Error(w, "Invalid Params", http.StatusBadRequest)
	}
}

func validAccessToken() bool {
	return false
}

func exchangeAuthCode(authCode string) (SonosAuthentication, error) {
	t := SonosAuthentication{}

	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("code", authCode)
	form.Add("redirect_uri", getRedirectURL())

	log.Info().
		Str("url", form.Encode()).
		Msg("encode form")

	call, err := makeAuthedRequest(form)
	if err != nil {
		return t, err
	}

	e := json.Unmarshal(call, &t)

	if e != nil {
		return t, e
	}

	return t, nil
}

func getCredentialsCode() (SonosAuthentication, error) {
	t := SonosAuthentication{}

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	call, err := makeAuthedRequest(form)
	if err != nil {
		return t, err
	}

	e := json.Unmarshal(call, &t)

	if e != nil {
		return t, e
	}

	return t, nil
}

func makeAuthedRequest(params url.Values) ([]byte, error) {
	// build request object
	ex, err := http.NewRequest("POST", fmt.Sprintf("%s/api/token", SPOTIFY_LOGIN), strings.NewReader(params.Encode()))
	ex.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ex.SetBasicAuth(config.ClientKey, config.ClientSecret)

	if err != nil {
		return nil, err
	}

	// execute request
	r, err := http.DefaultClient.Do(ex)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP call not successful (Status %d): %s", r.StatusCode, string(body))
	}

	return body, nil
}

func getRedirectURL() string {
	return fmt.Sprintf("http://localhost:%d/auth/redirect", config.Port)
}

func saveToken(token SonosAuthentication) error {
	content, e := json.MarshalIndent(token, "", "\t")

	if e != nil {
		return e
	}

	e = ioutil.WriteFile("token.json", content, 0644)
	if e != nil {
		return e
	}
	return nil
}
