package chipotle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kylegrantlucas/chipotle-go/menu"
	"github.com/kylegrantlucas/chipotle-go/search"
)

const BASE_URL = "https://services.chipotle.com"

type Client struct {
	APIKey     string
	httpClient *http.Client
}

// CustomTransport is a custom http.RoundTripper that adds default headers.
type customTransport struct {
	Transport http.RoundTripper
	Headers   map[string]string
}

// RoundTrip executes a single HTTP transaction and adds the custom headers.
func (c *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	newReq := req.Clone(req.Context())

	// Add the default headers
	for key, value := range c.Headers {
		newReq.Header.Set(key, value)
	}

	// Use the custom transport or fallback to http.DefaultTransport
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}

	return c.Transport.RoundTrip(newReq)
}

func NewClient(apiKey string) *Client {
	// create http client with api key header and content type set to application/json
	defaultHeaders := map[string]string{
		"Content-Type":              "application/json",
		"Ocp-Apim-Subscription-Key": apiKey,
	}

	// Create the custom transport
	customTransport := &customTransport{
		Headers: defaultHeaders,
	}

	// Create the http client
	// Use the custom transport
	client := &http.Client{
		Transport: customTransport,
	}

	return &Client{
		APIKey:     apiKey,
		httpClient: client,
	}
}

func (c *Client) Search(query search.Query) (*search.Result, error) {
	// marshal query to json
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	// create request
	req, err := http.NewRequest("POST", BASE_URL+"/restaurant/v3/restaurant", bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// print the body with the error
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("failed to execute request: %w, %s", err, string(body))
	}

	// check response status code
	if resp.StatusCode != http.StatusOK {
		// print the body with the error
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(body))
	}

	// decode response
	var result search.Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// handle pagination
	if result.PagingInfo.CurrentPage < result.PagingInfo.TotalPages {
		// create a new query with the next page
		// and recursively call Search
		// append the results to the current result
		// return the final result
		query.PageIndex = result.PagingInfo.CurrentPage + 1
		nextResult, err := c.Search(query)
		if err != nil {
			return nil, fmt.Errorf("failed to get next page: %w", err)
		}

		result.Restaurants = append(result.Restaurants, nextResult.Restaurants...)
	}

	return &result, nil
}

func (c *Client) GetMenu(restaurantID int) (*menu.Menu, error) {
	baseURL := "https://services.chipotle.com/menuinnovation/v1/restaurants/%d/onlinemenu?channelId=web&includeUnavailableItems=true"

	// create request
	req, err := http.NewRequest("GET", fmt.Sprintf(baseURL, restaurantID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// check response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// decode response
	var menu menu.Menu
	if err := json.NewDecoder(resp.Body).Decode(&menu); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &menu, nil
}
