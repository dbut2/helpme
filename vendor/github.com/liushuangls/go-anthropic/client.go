package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	config ClientConfig
}

// NewClient create new Anthropic API client
func NewClient(apikey string, opts ...ClientOption) *Client {
	return &Client{
		config: newConfig(apikey, opts...),
	}
}

func (c *Client) sendRequest(req *http.Request, v any) error {
	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := c.handlerRequestError(res); err != nil {
		return err
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) handlerRequestError(resp *http.Response) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&errRes)
		if err != nil || errRes.Error == nil {
			reqErr := RequestError{
				StatusCode: resp.StatusCode,
				Err:        err,
			}
			return &reqErr
		}
		return fmt.Errorf("error, status code: %d, message: %w", resp.StatusCode, errRes.Error)
	}
	return nil
}

func (c *Client) fullURL(suffix string) string {
	return fmt.Sprintf("%s%s", c.config.BaseURL, suffix)
}

func (c *Client) newRequest(ctx context.Context, method, urlSuffix string, body any) (req *http.Request, err error) {
	var reqBody []byte
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequestWithContext(ctx, method, c.fullURL(urlSuffix), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("X-Api-Key", c.config.apikey)
	req.Header.Set("Anthropic-Version", c.config.APIVersion)

	return req, nil
}

func (c *Client) newStreamRequest(ctx context.Context, method, urlSuffix string, body any) (req *http.Request, err error) {
	req, err = c.newRequest(ctx, method, urlSuffix, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	return req, nil
}
