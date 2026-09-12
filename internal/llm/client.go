package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bryantjandra/friday/internal/config"
)

/*
This Client is a wrapper around a generic http.Client. Our Client knows the api key, base url, and model,
whilst the inner http.Client knows how to send a request to a server and get a response back.
*/
type Client struct {
	apiKey  string
	baseUrl string
	model   string
	http    *http.Client
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		apiKey:  cfg.APIKey,
		baseUrl: cfg.BaseURL,
		model:   cfg.Model,
		http:    &http.Client{Timeout: 60 * time.Second},
	}
}

/* CreateMessage is a method, as we supply a receiver to it: (c *Client). This makes it a method for the Client type */
func (c *Client) CreateMessage(ctx context.Context, req Request) (*Response, error) {
	req.Model = c.model

	/* Takes the request struct and turns it into JSON (as an array of bytes). This JSON goes into the request body */
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	/* Build the url */
	url := c.baseUrl + "/v1/messages"

	/*
		We build the request.
		Arguments: 1 -> context so we can cancel the request, 2 -> it's a POST request, 3 -> where we are sending the request to, 4 -> the body (it needs a readable stream and not an array of raw bytes).
	*/
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	/* Send the request and wait for its response. Errors here will only be network type errors */
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}

	/* The response body is a stream that holds a network resource, and it must be closed. */
	/*
		What does it mean to hold a network resource?

		1. A network connection / socket was opened between my machine and the server (basically a channel over which bytes travel).
		2. My OS only allows a limited number of these open connections.

		If these open connections don't get closed, eventually we hit our OS' limit on open connections and new requests will then start failling.

		General principle: Anything that represents a live external resource (a network connection, a DB handle, a lock) must be explicitly closed or released when done using it.
	*/
	defer resp.Body.Close()

	/* Read the body --> The response body is a stream (data comes piece by piece over time), so io.ReadAll() reads the whole stream into an array of bytes */

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d, %s", resp.StatusCode, string(body))
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil

}
