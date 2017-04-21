package adzerk

import (
	"bytes"
	"io"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	defaultURL = "https://engine.adzerk.net/api/v2"
)

// A Client manages communication with the Adzerk API.
type Client struct {
	client *http.Client
	URL    string
}

// A RequestData is a parameter object which should be created and populated
// when generating a new Adzerk request using NewRequest.
type RequestData struct {
	IP               string
	UserID           string
	BlockedCreatives []int
	Keywords         []string
	Placements       []Placement
}

// A Request defines the JSON layout of the request body being sent to Adzerk.
type Request struct {
	IP               string      `json:"ip"`
	Time             int64       `json:"time"`
	Keywords         []string    `json:"keywords"`
	User             User        `json:"user"`
	Placements       []Placement `json:"placements"`
	BlockedCreatives []int       `json:"blockedCreatives"`
}

// A Placement defines the JSON layout of the `placements` object in the request body.
type Placement struct {
	NetworkID  int         `json:"networkId"`
	SiteID     int         `json:"siteId"`
	DivName    string      `json:"divName"`
	ZoneIDs    []int       `json:"zoneIds"`
	AdTypes    []int       `json:"adTypes"`
	Properties interface{} `json:"properties"`
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
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
	}

	return resp, err
}

// NewRequest creates a new POST API request.
func (c *Client) NewRequest(data RequestData) (*http.Request, error) {
	var buf io.ReadWriter
	body := Request{
		IP:       data.IP,
		Time:     time.Now().Unix(),
		Keywords: data.Keywords,
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
