// Package client provides a typed HTTP client for the Goban API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIResponse is the standard envelope returned by all Goban API endpoints.
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Client is an authenticated HTTP client for the Goban API.
type Client struct {
	BaseURL    string
	Token      string
	httpClient *http.Client
}

// New creates a new Client.
func New(baseURL, token string) *Client {
	return &Client{
		BaseURL:    baseURL,
		Token:      token,
		httpClient: &http.Client{},
	}
}

func (c *Client) do(method, path string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return c.httpClient.Do(req)
}

// decode reads and decodes an APIResponse[T] from resp.
func decode[T any](resp *http.Response) (T, error) {
	defer resp.Body.Close()
	var envelope APIResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		var zero T
		return zero, err
	}
	if !envelope.Success {
		var zero T
		msg := envelope.Error
		if msg == "" {
			msg = envelope.Message
		}
		return zero, fmt.Errorf("API error: %s", msg)
	}
	return envelope.Data, nil
}

func get[T any](c *Client, path string) (T, error) {
	resp, err := c.do(http.MethodGet, path, nil)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](resp)
}

func post[T any](c *Client, path string, body any) (T, error) {
	resp, err := c.do(http.MethodPost, path, body)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](resp)
}

func put[T any](c *Client, path string, body any) (T, error) {
	resp, err := c.do(http.MethodPut, path, body)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](resp)
}

func del[T any](c *Client, path string) (T, error) {
	resp, err := c.do(http.MethodDelete, path, nil)
	if err != nil {
		var zero T
		return zero, err
	}
	return decode[T](resp)
}
