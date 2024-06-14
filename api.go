package bdsaas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	apiPrefix = "https://openapi.bdsaas.com"
)

// Client is a client for the bdsaas API.
type Client struct {
	appKey string
}

// NewClient creates a new client with the given app key.
func NewClient(appKey string) *Client {
	return &Client{
		appKey: appKey,
	}
}

// GetSeats returns a list of phone numbers.
func (c Client) GetSeats(ctx context.Context) (seats []string, err error) {
	err = request(ctx, "POST", "/bdsaas/call/phoneApi/pagePhoneSeat.do", map[string]string{
		"appKey": c.appKey,
	}, &seats)
	return
}

// Call initiates a call from seat phone number to client phone number, with
// optional client IP address and notes, returns the session ID of the call.
func (c Client) Call(ctx context.Context, from, to string, options ...string) (sessionId string, err error) {
	var ip string
	var notes string
	if len(options) > 0 {
		ip = options[0]
		if len(options) > 1 {
			notes = options[1]
		}
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	err = request(ctx, "POST", "/bdsaas/call/phoneApi/callPhone.do", map[string]string{
		"appKey":    c.appKey,
		"seatPhone": from,
		"toPhone":   to,
		"ip":        ip,
		"ext1":      notes,
	}, &sessionId)
	return
}

// CallRecord represents a call record.
type CallRecord struct {
	CallerNum      string `json:"callerNum"`
	RecordFileID   string `json:"recordFileId"`
	DisconnectTime int64  `json:"disconnectTime"`
	TimeConsume    int    `json:"timeConsume"`
	SessionID      string `json:"sessionId"`
	CallType       string `json:"callType"`
	RealName       string `json:"realName"`
	CompanyID      int    `json:"companyId"`
	ProfileID      int    `json:"profileId"`
	ConnectTime    int64  `json:"connectTime"`
	CreatedTime    int64  `json:"createdTime"`
	CalleeNum      string `json:"calleeNum"`
	Status         string `json:"status"`
	Ext1           string `json:"ext1"`
}

// Query returns call records by session IDs.
func (c Client) Query(ctx context.Context, sessionIds ...string) (callRecords []CallRecord, err error) {
	err = request(ctx, "POST", "/bdsaas/call/phoneApi/queryCallPhoneRecord.do", map[string]interface{}{
		"appKey":     c.appKey,
		"sessionIds": sessionIds,
	}, &callRecords)
	return
}

type response struct {
	Code    int             `json:"rspCode"`
	Message string          `json:"rspMsg"`
	Data    json.RawMessage `json:"data"`
}

func request(ctx context.Context, method, path string, reqBody interface{}, target interface{}) error {
	var reader io.Reader = nil
	if method == "GET" {
		if queries, ok := reqBody.(map[string]string); ok {
			values := url.Values{}
			for k, v := range queries {
				values.Add(k, v)
			}
			path += "?" + values.Encode()
		}
	} else if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, apiPrefix+path, reader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	var resp response
	if err := json.Unmarshal(resBody, &resp); err != nil {
		var body string
		if len(resBody) > 1024 {
			body = string(resBody)[:1024] + "..."
		} else {
			body = string(resBody)
		}
		return fmt.Errorf("failed to decode response body (%s): %w", body, err)
	}
	if resp.Code != 0 {
		return fmt.Errorf("bdsaas api error: %s", resp.Message)
	}
	if err := json.Unmarshal(resp.Data, target); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}
	return nil
}
