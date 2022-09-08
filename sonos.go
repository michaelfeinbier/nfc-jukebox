package main

// Here we wrap (authenticated) sonos actions
type Sonos struct {
	AccessToken  string
	RefreshToken string
}

func New(accessToken string, refreshToken string) (s Sonos) {
	s = Sonos{AccessToken: accessToken, RefreshToken: refreshToken}
	return
}
