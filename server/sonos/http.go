package sonos

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type sonosHttp struct {
	client *http.Client
	host   string
}

const (
	SOAP_ENVELOPE = `<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <s:Body>%s</s:body>
	</s:envelope>`
)

func create(ip string) *sonosHttp {

	return &sonosHttp{
		host:   fmt.Sprintf("http://%s:1400", ip),
		client: &http.Client{},
	}
}

// We only do calls to the AVTransport service
func (h *sonosHttp) DoSopapCall(soapAction string, soapBody string) (string, error) {
	body := fmt.Sprintf(SOAP_ENVELOPE, soapBody)
	url := fmt.Sprintf("%s/MediaRenderer/AVTransport/Control", h.host)

	//log.Debug().
	//Str("url", url).
	//	Str("body", body).
	//	Msg("Send HTTP")

	//fmt.Println(body)

	req, e := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("SOAPAction", fmt.Sprintf("urn:schemas-upnp-org:service:AVTransport:1#%s", soapAction))
	if e != nil {
		return "", e
	}

	resp, e := h.client.Do(req)
	if e != nil {
		return "", e
	}

	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)

	return string(responseBody), nil
}
