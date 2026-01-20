// Package qi provides a Go client for the QiCard Payment Gateway API.
package qi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the default QiCard Payment Gateway API base URL.
	DefaultBaseURL = "https://api.qi.iq/api/v1"
	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Client is the QiCard Payment Gateway API client.
type Client struct {
	baseURL    string
	terminalID string
	username   string
	password   string
	signature  string
	httpClient *http.Client
}

// ClientOption is a function that configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithBasicAuth sets basic authentication credentials.
func WithBasicAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

// WithSignature sets the X-Signature header for signature-based authentication.
func WithSignature(signature string) ClientOption {
	return func(c *Client) {
		c.signature = signature
	}
}

// NewClient creates a new QiCard Payment Gateway API client.
func NewClient(terminalID string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    DefaultBaseURL,
		terminalID: terminalID,
		httpClient: &http.Client{Timeout: DefaultTimeout},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request and decodes the response.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Terminal-Id", c.terminalID)

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	if c.signature != "" {
		req.Header.Set("X-Signature", c.signature)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr Error
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Err:        &apiErr,
		}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// CreatePayment creates a new payment.
func (c *Client) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*Payment, error) {
	var payment Payment
	if err := c.doRequest(ctx, http.MethodPost, "/payment", req, &payment); err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetPaymentStatus retrieves the payment status by payment ID.
func (c *Client) GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentStatusResponse, error) {
	var status PaymentStatusResponse
	if err := c.doRequest(ctx, http.MethodGet, "/payment/"+paymentID+"/status", nil, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// GetPaymentStatusByRequest retrieves the payment status by request ID.
func (c *Client) GetPaymentStatusByRequest(ctx context.Context, requestID string) (*PaymentStatusResponse, error) {
	var status PaymentStatusResponse
	if err := c.doRequest(ctx, http.MethodGet, "/payment/status/by/request/"+requestID, nil, &status); err != nil {
		return nil, err
	}
	return &status, nil
}

// CancelPayment cancels a payment by payment ID.
func (c *Client) CancelPayment(ctx context.Context, paymentID string, req *CancelPaymentRequest) (*PaymentCancelResponse, error) {
	var resp PaymentCancelResponse
	if err := c.doRequest(ctx, http.MethodPost, "/payment/"+paymentID+"/cancel", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelPaymentByRequest cancels a payment by request ID.
func (c *Client) CancelPaymentByRequest(ctx context.Context, requestID string, req *CancelPaymentRequest) (*PaymentCancelResponse, error) {
	var resp PaymentCancelResponse
	if err := c.doRequest(ctx, http.MethodPost, "/payment/cancel/by/request/"+requestID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RefundPayment creates a refund for a payment by payment ID.
func (c *Client) RefundPayment(ctx context.Context, paymentID string, req *CreateRefundRequest) (*Refund, error) {
	var refund Refund
	if err := c.doRequest(ctx, http.MethodPost, "/payment/"+paymentID+"/refund", req, &refund); err != nil {
		return nil, err
	}
	return &refund, nil
}

// RefundPaymentByRequest creates a refund for a payment by request ID.
func (c *Client) RefundPaymentByRequest(ctx context.Context, requestID string, req *CreateRefundRequest) (*Refund, error) {
	var refund Refund
	if err := c.doRequest(ctx, http.MethodPost, "/payment/refund/by/request/"+requestID, req, &refund); err != nil {
		return nil, err
	}
	return &refund, nil
}
