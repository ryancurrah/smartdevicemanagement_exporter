package partnerconnmanager

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/api/smartdevicemanagement/v1"
)

var (
	authorizationURL = "https://nestservices.google.com/partnerconnections/%s/auth"
)

type AuthorizationCode struct {
	Code        string
	RedirectURI string
}

type PartnerConnManager struct {
	AuthorizationCodeChan chan AuthorizationCode
	ClientID              string
	ProjectID             string
}

func (p *PartnerConnManager) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(fmt.Sprintf(authorizationURL, p.ProjectID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	uri, err := redirectURI(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	v := url.Values{
		"access_type":   {"offline"},
		"prompt":        {"consent"},
		"response_type": {"code"},
		"redirect_uri":  {uri.String()},
		"client_id":     {p.ClientID},
		"scope":         {smartdevicemanagement.SdmServiceScope},
	}

	u.RawQuery = v.Encode()

	http.Redirect(w, r, u.String(), 301)
}

func (p *PartnerConnManager) AuthorizedHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	code := strings.TrimSpace(queries.Get("code"))
	if code == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "no authorization code received from partner connection manager: %s", r.URL.String())
		return
	}

	uri, err := redirectURI(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}

	p.AuthorizationCodeChan <- AuthorizationCode{Code: code, RedirectURI: uri.String()}

	fmt.Fprint(w, "authorization code received from partner connection manager")
}

func redirectURI(r *http.Request) (*url.URL, error) {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}

	u, err := url.Parse(fmt.Sprintf("%s://%s", scheme, r.Host))
	if err != nil {
		return nil, err
	}

	u.Path = r.URL.Path

	return u, nil
}
