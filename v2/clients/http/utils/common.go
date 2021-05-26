//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/salmanazmat666/go-mod-core-contracts/clients"
	"github.com/salmanazmat666/go-mod-core-contracts/errors"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos/common"

	"github.com/google/uuid"
)

// FromContext allows for the retrieval of the specified key's value from the supplied Context.
// If the value is not found, an empty string is returned.
func FromContext(ctx context.Context, key string) string {
	hdr, ok := ctx.Value(key).(string)
	if !ok {
		hdr = ""
	}
	return hdr
}

// correlatedId gets Correlation ID from supplied context. If no Correlation ID header is
// present in the supplied context, one will be created along with a value.
func correlatedId(ctx context.Context) string {
	correlation := FromContext(ctx, clients.CorrelationHeader)
	if len(correlation) == 0 {
		correlation = uuid.New().String()
	}
	return correlation
}

// Helper method to get the body from the response after making the request
func getBody(resp *http.Response) ([]byte, errors.EdgeX) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, errors.NewCommonEdgeX(errors.KindIOError, "failed to get the body from the response", err)
	}
	return body, nil
}

// Helper method to make the request and return the response
func makeRequest(req *http.Request) (*http.Response, errors.EdgeX) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, errors.NewCommonEdgeX(errors.KindClientError, "failed to send a http request", err)
	}
	return resp, nil
}

func createRequest(ctx context.Context, httpMethod string, baseUrl string, requestPath string, requestParams url.Values) (*http.Request, errors.EdgeX) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, "fail to parse baseUrl", err)
	}
	u.Path = requestPath
	if requestParams != nil {
		u.RawQuery = requestParams.Encode()
	}
	req, err := http.NewRequest(httpMethod, u.String(), nil)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, "failed to create a http request", err)
	}
	req.Header.Set(clients.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

func createRequestWithRawData(ctx context.Context, httpMethod string, url string, data interface{}) (*http.Request, errors.EdgeX) {
	jsonEncodedData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to encode input data to JSON", err)
	}

	content := FromContext(ctx, clients.ContentType)
	if content == "" {
		content = clients.ContentTypeJSON
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(jsonEncodedData))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, "failed to create a http request", err)
	}
	req.Header.Set(clients.ContentType, content)
	req.Header.Set(clients.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

// createRequestFromFilePath creates multipart/form-data request with the specified file
func createRequestFromFilePath(ctx context.Context, httpMethod string, url string, filePath string) (*http.Request, errors.EdgeX) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, fmt.Sprintf("fail to read file from %s", filePath), err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, "fail to create form data", err)
	}
	_, err = io.Copy(formFileWriter, bytes.NewReader(fileContents))
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, "fail to copy file to form data", err)
	}
	writer.Close()

	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindClientError, "failed to create a http request", err)
	}
	req.Header.Set(clients.ContentType, writer.FormDataContentType())
	req.Header.Set(clients.CorrelationHeader, correlatedId(ctx))
	return req, nil
}

// sendRequest will make a request with raw data to the specified URL.
// It returns the body as a byte array if successful and an error otherwise.
func sendRequest(ctx context.Context, req *http.Request) ([]byte, errors.EdgeX) {
	resp, err := makeRequest(req)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}
	if resp == nil {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, "the response should not be a nil", nil)
	}
	defer resp.Body.Close()

	bodyBytes, err := getBody(resp)
	if err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}

	if resp.StatusCode <= http.StatusMultiStatus {
		return bodyBytes, nil
	}

	// Handle error response
	var res common.BaseResponse
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return nil, errors.NewCommonEdgeXWrapper(err)
	}
	msg := fmt.Sprintf("request failed, status code: %d, err: %s", res.StatusCode, res.Message)
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, msg, nil)
	} else {
		return nil, errors.NewCommonEdgeX(errors.KindServerError, msg, nil)
	}
}
