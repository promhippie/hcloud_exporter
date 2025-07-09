package hetzner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/promhippie/hcloud_exporter/pkg/version"
)

var (
	// UserAgent defines the used user ganet for request.
	UserAgent = fmt.Sprintf("hcloud_exporter/%s", version.String)

	// Endpoint defines the endpoint for the Hetzner API.
	Endpoint = "https://api.hetzner.com/v1"
)

// Client is a client for the DockerHub API.
type Client struct {
	httpClient *http.Client
	token      string
	StorageBox StorageBoxClient
}

// A ClientOption is used to configure a Client.
type ClientOption func(*Client)

// WithToken configures a Client to use the specified token for authentication.
func WithToken(value string) ClientOption {
	return func(client *Client) {
		client.token = value
	}
}

// NewClient creates a new client.
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	client.StorageBox = StorageBoxClient{
		client: client,
	}

	return client
}

// NewRequest creates an HTTP request against the DockerHub.
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	if c.token != "" {
		req.Header.Set(
			"Authorization",
			fmt.Sprintf(
				"Bearer %s",
				c.token,
			),
		)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req.WithContext(ctx), nil
}

// Do performs an HTTP request against the DockerHub.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	res, err := c.httpClient.Do(req)

	if res != nil {
		defer func() { _ = res.Body.Close() }()
	}

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return &Response{Response: res}, err
	}

	res.Body = io.NopCloser(bytes.NewReader(body))

	if res.StatusCode >= 400 && res.StatusCode <= 599 {
		return &Response{Response: res}, errors.New(http.StatusText(res.StatusCode))
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, bytes.NewReader(body))
		} else {
			err = json.Unmarshal(body, v)
		}
	}

	return &Response{Response: res}, err
}

// Response simply wraps the standard response type.
type Response struct {
	*http.Response
}
