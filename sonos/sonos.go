package sonos

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

const (
	SPOTIFY_METADATA = `
<DIDL-Lite xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:upnp="urn:schemas-upnp-org:metadata-1-0/upnp/" xmlns:r="urn:schemas-rinconnetworks-com:metadata-1-0/" xmlns="urn:schemas-upnp-org:metadata-1-0/DIDL-Lite/">
    <item id="%s" restricted="true">
        <dc:title></dc:title>
        <upnp:class>%s</upnp:class>
        <desc id="cdudn" nameSpace="urn:schemas-rinconnetworks-com:metadata-1-0/">SA_RINCON2311_X_#Svc2311-0-Token</desc>
    </item>
</DIDL-Lite>
`
	UPNP_CLASS_TRACK = "object.item.audioItem.musicTrack"
	UPNP_CLASS_ALBUM = "object.container.album.musicAlbum"

	// TODO: put to config?
	COORDINATOR_UID = "RINCON_7828CA0E099601400"
)

type SonosPlayer struct {
	ip   string
	port int8
}

type soapRemoveAllTracksFromQueue struct {
	XMLName xml.Name `xml:"u:RemoveAllTracksFromQueue"`
	XmlNS   string   `xml:"xmlns:test,attr"`
}

var httpHelper *sonosHttp

func New(ipAddress string) (s SonosPlayer) {
	s = SonosPlayer{ip: ipAddress}
	httpHelper = create(ipAddress)
	return
}

// Plays the item on top of the current queue
func (s *SonosPlayer) Play() error {
	body := `<u:Play xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
            <InstanceID>0</InstanceID>
            <Speed>1</Speed>
        </u:Play>`

	_, e := httpHelper.DoSopapCall("Play", body)
	return e
}

// Removes all Tracks from the current queue
func (s *SonosPlayer) RemoveAllTracksFromQueue() error {
	body := `<u:RemoveAllTracksFromQueue xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
            <InstanceID>0</InstanceID>
        </u:RemoveAllTracksFromQueue>`

	_, e := httpHelper.DoSopapCall("RemoveAllTracksFromQueue", body)
	return e
}

// Adds an album to the current queue
func (s *SonosPlayer) AddSpotifyAlbumToQueue(spotifyURI string) error {
	spotifyURIEncoded, spotifyMetadata := getMetadataForAlbum(spotifyURI)
	body := fmt.Sprintf(`<u:AddURIToQueue xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
            <InstanceID>0</InstanceID>
            <EnqueuedURI>%s</EnqueuedURI>
            <EnqueuedURIMetaData>%s</EnqueuedURIMetaData>
            <DesiredFirstTrackNumberEnqueued>0</DesiredFirstTrackNumberEnqueued>
            <EnqueueAsNext>1</EnqueueAsNext>
        </u:AddURIToQueue>`, spotifyURIEncoded, spotifyMetadata)

	_, e := httpHelper.DoSopapCall("AddURIToQueue", body)
	return e
}

// set the input to the queue
func (s *SonosPlayer) SetInputToQueue() error {
	body := fmt.Sprintf(`<u:SetAVTransportURI xmlns:u="urn:schemas-upnp-org:service:AVTransport:1">
  <InstanceID>0</InstanceID>
  <CurrentURI>x-rincon-queue:%s#0</CurrentURI>
  <CurrentURIMetaData></CurrentURIMetaData>
</u:SetAVTransportURI>`, COORDINATOR_UID)
	_, e := httpHelper.DoSopapCall("SetAVTransportURI", body)
	return e
}

func (s *SonosPlayer) PlaySpotifyAlbum(spotifyURI string) error {
	err := s.SetInputToQueue()
	if err != nil {
		return err
	}

	err = s.RemoveAllTracksFromQueue()
	if err != nil {
		return err
	}

	err = s.AddSpotifyAlbumToQueue(spotifyURI)
	if err != nil {
		return err
	}

	err = s.Play()

	return err
}

func getMetadataForAlbum(uri string) (trackUri, body string) {
	trackUri = "x-rincon-cpcontainer:1004206c" + url.QueryEscape(uri) + "?sid=9&amp;flags=8300&amp;sn=7"
	itemid := "0004206c" + url.QueryEscape(uri)
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;", "\"", "&quot;")
	body = r.Replace(fmt.Sprintf(SPOTIFY_METADATA, itemid, UPNP_CLASS_ALBUM))

	return
}
