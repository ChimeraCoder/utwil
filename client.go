package utwil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// At the time of writing, the current API version was released on Apr. 1, 2010
const (
	BaseURL   = "https://api.twilio.com"
	LookupURL = "https://lookups.twilio.com/v1"

	APIVersion = "2010-04-01"
)

// Client stores Twilio API credentials
type Client struct {
	AccountSID string
	AuthToken  string
}

func (c *Client) getJSON(url string, result interface{}) error {
	hc := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("GetJSON(): %s", err)
	}
	req.SetBasicAuth(c.AccountSID, c.AuthToken)
	resp, err := hc.Do(req)
	if err != nil {
		return fmt.Errorf("GetJSON(): %s", err)
	}

	if resp.StatusCode != 200 {
		re := RESTException{}
		json.NewDecoder(resp.Body).Decode(&re)
		return re
	}
	return json.NewDecoder(resp.Body).Decode(&result)
}

func (c *Client) postForm(url string, values url.Values, result interface{}) error {
	hc := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("PostForm(): %s", err)
	}
	req.SetBasicAuth(c.AccountSID, c.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err := RESTException{}
		json.NewDecoder(resp.Body).Decode(&err)
		return err
	}
	return json.NewDecoder(resp.Body).Decode(&result)
}

func (c *Client) urlPrefix() string {
	return fmt.Sprintf("%s/%s/Accounts/%s", BaseURL, APIVersion, c.AccountSID)
}

func (c *Client) callsURL() string {
	return fmt.Sprintf("%s/Calls.json", c.urlPrefix())
}

func (c *Client) messagesURL() string {
	return fmt.Sprintf("%s/Messages.json", c.urlPrefix())
}
