package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type BaseClient struct {
	BaseURL     string
	Client      *http.Client
	Headers     Headers
	QueryParams Query
}

type Query map[string]string
type Headers map[string]string

func (c *BaseClient) ExecuteRequestAndGetResponse(ctx context.Context, method string, url string, queryParams Query, headers Headers, data interface{}, result interface{}) error {
	responseBody, err := c.ExecuteRequest(ctx, method, url, queryParams, headers, data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseBody, result)
	if err != nil {
		return err
	}

	return nil
}

func (c *BaseClient) ExecuteRequest(ctx context.Context, method string, requestURL string, queryParams Query, headers Headers, data interface{}) ([]byte, error) {
	var body io.Reader

	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	query := req.URL.Query()
	for k, v := range c.QueryParams {
		query.Add(k, v)
	}
	for k, v := range queryParams {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()

	rsp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"method":     method,
		"requestURL": req.URL.String(),
		"body":       data,
		"response":   string(rspBody),
	}).Debug("HTTP CLIENT REQUEST")

	return rspBody, nil
}

func (c *BaseClient) BuildURL(url string, args ...interface{}) string {
	urlTemplate := fmt.Sprintf("%s%s", c.BaseURL, url)
	return fmt.Sprintf(urlTemplate, args...)
}
