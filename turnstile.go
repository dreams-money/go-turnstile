package turnstile

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

var ErrTimeoutOrDuplicate = errors.New("timeout or duplicate request")

type SiteVerifyResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
	Messages   []string `json:"messages"`
	// Hostname   string   `json:"hostname"`
	// TokenID    string   `json:"tokenId"`
}

type Client struct {
	SecretKey string
}

func New(secretKey string) Client {
	return Client{
		SecretKey: secretKey,
	}
}

func (c *Client) Verify(token, remoteip string) (error, error) {

	siteVerifyForm := url.Values{}
	siteVerifyForm.Add("secret", c.SecretKey)
	siteVerifyForm.Add("response", token)
	siteVerifyForm.Add("remoteip", remoteip)
	siteVerifyFormData := strings.NewReader(siteVerifyForm.Encode())

	siteVerifyURL := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	siteVerifyReq, err := http.NewRequest(http.MethodPost, siteVerifyURL, siteVerifyFormData)
	if err != nil {
		return nil, err
	}

	siteVerifyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	siteVerifyResp, err := (&http.Client{}).Do(siteVerifyReq)
	if err != nil {
		return nil, err
	}
	defer siteVerifyResp.Body.Close()

	siteVerifyRespBody, err := io.ReadAll(siteVerifyResp.Body)
	if err != nil {
		return nil, err
	}

	siteVerifyResponse := SiteVerifyResponse{}
	err = json.Unmarshal(siteVerifyRespBody, &siteVerifyResponse)
	if err != nil {
		return nil, err
	}

	if !siteVerifyResponse.Success {
		return siteVerifyResponse.error(), nil
	}

	return nil, nil
}

func (r SiteVerifyResponse) error() error {
	if slices.Contains(r.ErrorCodes, "timeout-or-duplicate") {
		return ErrTimeoutOrDuplicate
	}

	return errors.New(r.errorsToString())
}

func (r SiteVerifyResponse) errorsToString() string {
	codes := ""
	for _, code := range r.ErrorCodes {
		codes += code + ","
	}
	codes = strings.TrimRight(codes, ",")
	s := "[" + codes + "]"

	if len(r.Messages) > 0 {
		messages := ""
		for i, message := range r.Messages {
			messageNumber := strconv.Itoa(i + 1)
			messages += messageNumber + ": " + message + ", "
		}
		messages = strings.TrimRight(messages, ", ")
		s += "; Messages: " + messages
	}

	return s
}
