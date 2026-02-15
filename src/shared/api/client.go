package api

import (
	"strconv"
	"time"

	"resty.dev/v3"
)

type APIClient struct {
	client  *resty.Client
	baseURL string
	csot    string
}

func NewAPIClient(baseURL, csot string) *APIClient {
	client := resty.New().
		SetBaseURL(baseURL).
		SetTimeout(15*time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(500*time.Millisecond).
		SetRetryMaxWaitTime(2*time.Second).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("x-internal-token", csot)

	return &APIClient{
		client:  client,
		baseURL: baseURL,
		csot:    csot,
	}
}

func setUserHeaders(r *resty.Request, userID int64, username string) {
	r.SetHeader("x-internal-user-id", strconv.FormatInt(userID, 10))
	if username != "" {
		r.SetHeader("x-internal-user-name", username)
	}
}

func validateResponse(resp *resty.Response) error {
	if resp.IsError() {
		return &APIError{
			StatusCode: resp.StatusCode(),
			Message:    resp.String(),
		}
	}
	return nil
}
