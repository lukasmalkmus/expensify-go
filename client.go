package expensify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL = "https://integrations.expensify.com/Integration-Server/ExpensifyIntegrations"
	version = "0.1.0"
)

var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

// Error is the generic error response returned on non 2xx HTTP status codes.
type Error struct {
	Message    string `json:"responseMessage"`
	StatusCode int    `json:"responseCode"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(e.StatusCode), e.Message)
}

// An Option can be used to configure the behaviour of the API client.
type Option func(c *Client) error

// SetClient specifies a custom http client that should be used to make
// requests.
func SetClient(client *http.Client) Option {
	return func(c *Client) error {
		if client == nil {
			return nil
		}
		c.httpClient = client
		return nil
	}
}

// Client provides the Expensify HTTP API operations.
type Client struct {
	baseURL   *url.URL
	userAgent string

	partnerUserID     string
	partnerUserSecret string

	httpClient *http.Client

	Expense ExpenseService
}

// NewClient returns a new Expensify API client. The credentials can be
// retrieved from https://www.expensify.com/tools/integrations.
func NewClient(partnerUserID, partnerUserSecret string, options ...Option) (*Client, error) {
	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL:   u,
		userAgent: fmt.Sprintf("expensify-go/%s", version),

		partnerUserID:     partnerUserID,
		partnerUserSecret: partnerUserSecret,

		httpClient: defaultHTTPClient,
	}
	c.Expense = &expenseService{c}

	// Apply supplied options.
	if err := c.Options(options...); err != nil {
		return nil, err
	}

	return c, nil
}

// Options applies Options to a client instance.
func (c *Client) Options(options ...Option) error {
	for _, option := range options {
		if err := option(c); err != nil {
			return err
		}
	}
	return nil
}

// call creates a new API request and executes it.
func (c *Client) call(ctx context.Context, jobType, inputType string, payload, v interface{}) error {
	req, err := c.newRequest(ctx, jobType, inputType, payload)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// newRequest creates an API request. If specified, the value pointed to by
// body will be included as the request body.
func (c *Client) newRequest(ctx context.Context, jobType, inputType string, payload interface{}) (*http.Request, error) {
	job := &jobRequest{
		Type: jobType,
		InputSettings: &inputSettings{
			Type: inputType,
			data: payload,
		},
	}
	job.Credentials.PartnerUserID = c.partnerUserID
	job.Credentials.PartnerUserSecret = c.partnerUserSecret

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(job); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("requestJobDescription", buf.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	// Set headers.
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

// do sends an API request and returns the API response.
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// The Expensify API is super shitty and doesn't return proper HTTP status
	// codes (200 all the way). So we need to decode the body and see if it's an
	// error. If not, we need to decode it again into the proper response
	// struct.

	var (
		buf bytes.Buffer
		r   = io.TeeReader(resp.Body, &buf)
	)

	var errResp Error
	if err = json.NewDecoder(r).Decode(&errResp); err != nil {
		return err
	} else if code := errResp.StatusCode; code != 0 && code != http.StatusOK {
		return errResp
	}

	if v != nil {
		if err = json.NewDecoder(&buf).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
