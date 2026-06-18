package payment

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const snapSandboxURL = "https://app.sandbox.midtrans.com/snap/v1/transactions"

type Client struct {
	serverKey string
	clientKey string
	httpc     *http.Client
}

func NewClient() *Client {
	return &Client{
		serverKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		clientKey: os.Getenv("MIDTRANS_CLIENT_KEY"),
		httpc:     &http.Client{},
	}
}

func (c *Client) ClientKey() string { return c.clientKey }

type SnapRequest struct {
	OrderID     string
	Amount      int
	CustomerName string
	Email       string
}

type snapResponse struct {
	Token         string `json:"token"`
	RedirectURL   string `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages"`
}

// CreateTransaction calls Midtrans Snap API and returns the Snap token.
func (c *Client) CreateTransaction(req SnapRequest) (string, error) {
	body := map[string]any{
		"transaction_details": map[string]any{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
		"customer_details": map[string]any{
			"first_name": req.CustomerName,
			"email":      req.Email,
		},
		"item_details": []map[string]any{
			{
				"id":       "stnk-renewal",
				"price":    req.Amount,
				"quantity": 1,
				"name":     "Perpanjangan STNK",
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal snap request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, snapSandboxURL, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	// Basic auth: serverKey as username, empty password
	auth := base64.StdEncoding.EncodeToString([]byte(c.serverKey + ":"))
	httpReq.Header.Set("Authorization", "Basic "+auth)

	resp, err := c.httpc.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("snap request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var sr snapResponse
	if err := json.Unmarshal(respBody, &sr); err != nil {
		return "", fmt.Errorf("unmarshal snap response: %w (body: %s)", err, string(respBody))
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("snap error (status %d): %v", resp.StatusCode, sr.ErrorMessages)
	}

	if sr.Token == "" {
		return "", fmt.Errorf("snap returned empty token: %s", string(respBody))
	}

	return sr.Token, nil
}

// Notification is the payload Midtrans sends to the webhook.
type Notification struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

// VerifySignature checks the SHA512 signature to ensure the notification is genuine.
// Formula: sha512(order_id + status_code + gross_amount + server_key)
func (c *Client) VerifySignature(n Notification) bool {
	raw := n.OrderID + n.StatusCode + n.GrossAmount + c.serverKey
	sum := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(sum[:])
	return expected == n.SignatureKey
}

// IsSuccess reports whether the transaction status means a completed payment.
func (n Notification) IsSuccess() bool {
	switch n.TransactionStatus {
	case "capture":
		return n.FraudStatus == "accept" || n.FraudStatus == ""
	case "settlement":
		return true
	default:
		return false
	}
}

// IsFailure reports whether the transaction failed/expired/cancelled.
func (n Notification) IsFailure() bool {
	switch n.TransactionStatus {
	case "deny", "cancel", "expire", "failure":
		return true
	default:
		return false
	}
}
