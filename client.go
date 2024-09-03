package propel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

)

type Client struct {
	BaseURL string
	HTTPClient *http.Client
	APIKey string
}

func NewClient(base_url string, apiKey string) *Client {
    return &Client{
        BaseURL: base_url,
        HTTPClient: &http.Client{},
        APIKey: apiKey,
    }
}

func (c *Client) makeRequest(method string, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	var requestBody []byte
    var err error

    if body != nil {
        requestBody, err = json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal body: %w", err)
        }
    }

    url := fmt.Sprintf("%s%s", c.BaseURL, path)
    request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    for key, value := range headers {
        request.Header.Set(key, value)
    }

    request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
    request.Header.Set("Content-Type", "application/json")

    response, err := c.HTTPClient.Do(request)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }

    if response.StatusCode < 200 || response.StatusCode >= 300 {
        var buf bytes.Buffer
        _, err := io.Copy(&buf, response.Body)
        if err != nil {
            return nil, fmt.Errorf("failed to read response body: %w", err)
        }
        return response, fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode, buf.String())
    }

    return response, nil
}

func (c *Client) GetUserByID(id string) (*User, error) {
	path := fmt.Sprintf("/users/%s", id)

    response, err := c.makeRequest("GET", path, nil, nil) 
    if err != nil {
        return nil, err
    }

    defer response.Body.Close()

    var user User
    if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
        return nil, err
    }

    return &user, nil
}