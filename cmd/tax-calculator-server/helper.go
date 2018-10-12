package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func httpGet(appUrl string, authToken string, params *url.Values) (map[string]interface{}, int, error) {
	if params == nil {
		params = &url.Values{}
	}

	var fullUrl string
	if getParams := params.Encode(); getParams == "" {
		fullUrl = appUrl
	} else {
		fullUrl = fmt.Sprintf("%s?%s", appUrl, getParams)
	}

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		return nil, -1, err
	}

	req.Header.Set("Authentication-Token", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, -1, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(respBody, &payload)
	if err != nil {
		return nil, -1, err
	}

	return payload, resp.StatusCode, nil
}

func httpPost(endpoint string, authToken string, params *url.Values) (map[string]interface{}, int, error) {
	if params == nil {
		params = &url.Values{}
	}

	contentType := "application/x-www-form-urlencoded"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, -1, err
	}
	req.Header.Set("content-type", contentType)
	req.Header.Set("Authentication-Token", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, -1, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(respBody, &payload)
	if err != nil {
		return nil, -1, err
	}

	return payload, resp.StatusCode, nil
}
