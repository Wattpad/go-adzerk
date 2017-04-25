package adzerk

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	defaultURL = "https://engine.adzerk.net/api/v2"
)

var (
	ErrMissingAdTypes   = errors.New("go-adzerk/adzerk: missing value for Placements.AdTypes")
	ErrMissingDivName   = errors.New("go-adzerk/adzerk: missing value for Placements.NetworkID")
	ErrMissingNetworkID = errors.New("go-adzerk/adzerk: missing value for Placements.NetworkID")
	ErrMissingSiteID    = errors.New("go-adzerk/adzerk: missing value for Placements.SiteID")
	ErrNoPlacements     = errors.New("go-adzerk/adzerk: missing value for RequestData.Placements")
)

// A Client manages communication with the Adzerk API.
type Client struct {
	client *http.Client
	URL    string
}

// A RequestData is a parameter object which should be created and populated
// when generating a new Adzerk request using NewRequest.
type RequestData struct {
	// Required arguments
	Placements []Placement

	// Optional arguments
	EnableBotFiltering bool
	IncludePricingData bool
	IsMobile           bool
	NoTrack            bool
	BlockedCreatives   []int
	IP                 string
	Referrer           string
	URL                string
	UserID             string
	Keywords           []string
}

// A Request defines the JSON layout of the request body being sent to Adzerk.
type Request struct {
	EnableBotFiltering bool        `json:"enableBotFiltering,omitempty"`
	IncludePricingData bool        `json:"includePricingData,omitempty"`
	IsMobile           bool        `json:"isMobile,omitempty"`
	NoTrack            bool        `json:"notrack,omitempty"`
	Time               int64       `json:"time,omitempty"`
	BlockedCreatives   []int       `json:"blockedCreatives,omitempty"`
	IP                 string      `json:"ip,omitempty"`
	Referrer           string      `json:"referrer,omitempty"`
	URL                string      `json:"url,omitempty"`
	Keywords           []string    `json:"keywords,omitempty"`
	User               User        `json:"user,omitempty"`
	Placements         []Placement `json:"placements"`
}

// A Placement defines the JSON layout of the `placements` object in the request body.
type Placement struct {
	// Required arguments
	NetworkID int    `json:"networkId"`
	SiteID    int    `json:"siteId"`
	AdTypes   []int  `json:"adTypes"`
	DivName   string `json:"divName"`

	// Optional arguments
	AdID        int         `json:"adId,omitempty"`
	CampaignID  int         `json:"campaignId,omitempty"`
	FlightID    int         `json:"flightId,omitempty"`
	EventIDs    []int       `json:"eventIds,omitempty"`
	ZoneIDs     []int       `json:"zoneIds,omitempty"`
	ClickURL    string      `json:"clickUrl,omitempty"`
	ContentKeys interface{} `json:"contentKeys,omitempty"`
	Overrides   interface{} `json:"overrides,omitempty"`
	Properties  interface{} `json:"properties,omitempty"`
}

// A User defines the JSON layout of the `user` object in the request body.
type User struct {
	Key string `json:"key"`
}

// NewClient returns a new Adzerk client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = http.DefaultClient
	}

	return &Client{
		client: c,
		URL:    defaultURL,
	}
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := ctxhttp.Do(ctx, c.client, req)
	if err != nil {
		return nil, errors.Wrap(err, "go-adzerk/adzerk")
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "go-adzerk/adzerk")
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return nil, errors.Wrap(err, "go-adzerk/adzerk")
	}

	return resp, err
}

// NewRequest creates a new POST API request.
func (c *Client) NewRequest(data RequestData) (*http.Request, error) {
	if len(data.Placements) == 0 {
		return nil, ErrNoPlacements
	}
	for _, p := range data.Placements {
		if p.NetworkID == 0 {
			return nil, ErrMissingNetworkID
		}
		if p.SiteID == 0 {
			return nil, ErrMissingSiteID
		}
		if len(p.AdTypes) == 0 {
			return nil, ErrMissingAdTypes
		}
		if p.DivName == "" {
			return nil, ErrMissingDivName
		}
	}
	var buf io.ReadWriter
	body := Request{
		EnableBotFiltering: data.EnableBotFiltering,
		IncludePricingData: data.IncludePricingData,
		IsMobile:           data.IsMobile,
		NoTrack:            data.NoTrack,
		IP:                 data.IP,
		Referrer:           data.Referrer,
		URL:                data.URL,
		Time:               time.Now().Unix(),
		Keywords:           data.Keywords,
		User: User{
			Key: data.UserID,
		},
		Placements:       data.Placements,
		BlockedCreatives: data.BlockedCreatives,
	}

	buf = new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.URL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
