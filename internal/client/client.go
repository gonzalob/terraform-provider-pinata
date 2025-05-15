package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const HostURL string = "https://api.pinata.cloud"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(host, token *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if token == nil {
		return &c, nil
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Token

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
