NFC-Vinyl Player
==================

This is a small hobby-project of mine to learn Golang and to help automating my smart home.

## Basic Idea
Idea is to put NFC Tags on all the vinyls I own. If you scan a tag with your phone it will play the album on Spotify/Sonos.

I plan to achieve this by running a small webserver on a raspberry pi and then put a unique URI per vinyl on a NFC Tag. Once this URI is called the playback could start immediately.

## What we need
* [x] Webserver
* [ ] Authentication against Spotify/Sonos/Apple Music (tbd.)
* [ ] Databaselike key/value store for NFC Tag(URL?) to Album matching (maybe [Bitcask?](https://pkg.go.dev/git.mills.io/prologic/bitcask))
* [ ] A lot of NFC Tags :wink:

## Open Questions
* How do we automatically play a specific Album (Spotify URI) on a specific Sonos zone player?

## Links / Inspiration
- https://www.instructables.com/Albums-With-NFC-Tags-to-Automatically-Play-Spotify/
- https://shkspr.mobi/blog/2020/09/how-can-i-launch-a-spotify-album-from-an-nfc-tag/
- https://developer.sonos.com/reference/control-api/

https://open.spotify.com/track/2SgQjReOSmKZYnRnt88B09?si=7f0e16d0bca9460a
