// Package client provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// SchemaObject defines model for SchemaObject.
type SchemaObject struct {
	FirstName string `json:"firstName"`
	Role      string `json:"role"`
}

// PostBothJSONBody defines parameters for PostBoth.
type PostBothJSONBody SchemaObject

// PostJsonJSONBody defines parameters for PostJson.
type PostJsonJSONBody SchemaObject

// PostBothRequestBody defines body for PostBoth for application/json ContentType.
type PostBothJSONRequestBody PostBothJSONBody

// PostJsonRequestBody defines body for PostJson for application/json ContentType.
type PostJsonJSONRequestBody PostJsonJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

var _ ClientInterface = &Client{}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// NewClient Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostBothWithBody request  with any body
	PostBothWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)
	// PostBothWithBodyWithResponse request  with any body and parse response
	PostBothWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostBothResponse, error)

	// PostBoth
	PostBoth(ctx context.Context, body PostBothJSONRequestBody) (*http.Response, error)
	// PostBothWithResponse
	PostBothWithResponse(ctx context.Context, body PostBothJSONRequestBody) (*PostBothResponse, error)

	// GetBoth request
	GetBoth(ctx context.Context) (*http.Response, error)
	// GetBothWithResponse request  and parse response
	GetBothWithResponse(ctx context.Context) (*GetBothResponse, error)

	// PostJsonWithBody request  with any body
	PostJsonWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)
	// PostJsonWithBodyWithResponse request  with any body and parse response
	PostJsonWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostJsonResponse, error)

	// PostJson
	PostJson(ctx context.Context, body PostJsonJSONRequestBody) (*http.Response, error)
	// PostJsonWithResponse
	PostJsonWithResponse(ctx context.Context, body PostJsonJSONRequestBody) (*PostJsonResponse, error)

	// GetJson request
	GetJson(ctx context.Context) (*http.Response, error)
	// GetJsonWithResponse request  and parse response
	GetJsonWithResponse(ctx context.Context) (*GetJsonResponse, error)

	// PostOtherWithBody request  with any body
	PostOtherWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)
	// PostOtherWithBodyWithResponse request  with any body and parse response
	PostOtherWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostOtherResponse, error)

	// GetOther request
	GetOther(ctx context.Context) (*http.Response, error)
	// GetOtherWithResponse request  and parse response
	GetOtherWithResponse(ctx context.Context) (*GetOtherResponse, error)

	// GetJsonWithTrailingSlash request
	GetJsonWithTrailingSlash(ctx context.Context) (*http.Response, error)
	// GetJsonWithTrailingSlashWithResponse request  and parse response
	GetJsonWithTrailingSlashWithResponse(ctx context.Context) (*GetJsonWithTrailingSlashResponse, error)
}

func (c *Client) PostBothWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewPostBothRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) PostBoth(ctx context.Context, body PostBothJSONRequestBody) (*http.Response, error) {
	req, err := NewPostBothRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetBoth(ctx context.Context) (*http.Response, error) {
	req, err := NewGetBothRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) PostJsonWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewPostJsonRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) PostJson(ctx context.Context, body PostJsonJSONRequestBody) (*http.Response, error) {
	req, err := NewPostJsonRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetJson(ctx context.Context) (*http.Response, error) {
	req, err := NewGetJsonRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) PostOtherWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewPostOtherRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetOther(ctx context.Context) (*http.Response, error) {
	req, err := NewGetOtherRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetJsonWithTrailingSlash(ctx context.Context) (*http.Response, error) {
	req, err := NewGetJsonWithTrailingSlashRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewPostBothRequest calls the generic PostBoth builder with application/json body
func NewPostBothRequest(server string, body PostBothJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostBothRequestWithBody(server, "application/json", bodyReader)
}

// NewPostBothRequestWithBody generates requests for PostBoth with any type of body
func NewPostBothRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_both_bodies")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewGetBothRequest generates requests for GetBoth
func NewGetBothRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_both_responses")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostJsonRequest calls the generic PostJson builder with application/json body
func NewPostJsonRequest(server string, body PostJsonJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostJsonRequestWithBody(server, "application/json", bodyReader)
}

// NewPostJsonRequestWithBody generates requests for PostJson with any type of body
func NewPostJsonRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_json_body")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewGetJsonRequest generates requests for GetJson
func NewGetJsonRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_json_response")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostOtherRequestWithBody generates requests for PostOther with any type of body
func NewPostOtherRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_other_body")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewGetOtherRequest generates requests for GetOther
func NewGetOtherRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_other_response")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetJsonWithTrailingSlashRequest generates requests for GetJsonWithTrailingSlash
func NewGetJsonWithTrailingSlashRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/with_trailing_slash/")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

type PostBothResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r PostBothResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostBothResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetBothResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SchemaObject
}

// Status returns HTTPResponse.Status
func (r GetBothResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetBothResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostJsonResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r PostJsonResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostJsonResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetJsonResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SchemaObject
}

// Status returns HTTPResponse.Status
func (r GetJsonResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetJsonResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostOtherResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r PostOtherResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostOtherResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetOtherResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetOtherResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetOtherResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetJsonWithTrailingSlashResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *SchemaObject
}

// Status returns HTTPResponse.Status
func (r GetJsonWithTrailingSlashResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetJsonWithTrailingSlashResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostBothWithBodyWithResponse request with arbitrary body returning *PostBothResponse
func (c *Client) PostBothWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostBothResponse, error) {
	rsp, err := c.PostBothWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsePostBothResponse(rsp)
}

func (c *Client) PostBothWithResponse(ctx context.Context, body PostBothJSONRequestBody) (*PostBothResponse, error) {
	rsp, err := c.PostBoth(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsePostBothResponse(rsp)
}

// GetBothWithResponse request returning *GetBothResponse
func (c *Client) GetBothWithResponse(ctx context.Context) (*GetBothResponse, error) {
	rsp, err := c.GetBoth(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetBothResponse(rsp)
}

// PostJsonWithBodyWithResponse request with arbitrary body returning *PostJsonResponse
func (c *Client) PostJsonWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostJsonResponse, error) {
	rsp, err := c.PostJsonWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsePostJsonResponse(rsp)
}

func (c *Client) PostJsonWithResponse(ctx context.Context, body PostJsonJSONRequestBody) (*PostJsonResponse, error) {
	rsp, err := c.PostJson(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsePostJsonResponse(rsp)
}

// GetJsonWithResponse request returning *GetJsonResponse
func (c *Client) GetJsonWithResponse(ctx context.Context) (*GetJsonResponse, error) {
	rsp, err := c.GetJson(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetJsonResponse(rsp)
}

// PostOtherWithBodyWithResponse request with arbitrary body returning *PostOtherResponse
func (c *Client) PostOtherWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostOtherResponse, error) {
	rsp, err := c.PostOtherWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsePostOtherResponse(rsp)
}

// GetOtherWithResponse request returning *GetOtherResponse
func (c *Client) GetOtherWithResponse(ctx context.Context) (*GetOtherResponse, error) {
	rsp, err := c.GetOther(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetOtherResponse(rsp)
}

// GetJsonWithTrailingSlashWithResponse request returning *GetJsonWithTrailingSlashResponse
func (c *Client) GetJsonWithTrailingSlashWithResponse(ctx context.Context) (*GetJsonWithTrailingSlashResponse, error) {
	rsp, err := c.GetJsonWithTrailingSlash(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetJsonWithTrailingSlashResponse(rsp)
}

// ParsePostBothResponse parses an HTTP response from a PostBothWithResponse call
func ParsePostBothResponse(rsp *http.Response) (*PostBothResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &PostBothResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetBothResponse parses an HTTP response from a GetBothWithResponse call
func ParseGetBothResponse(rsp *http.Response) (*GetBothResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetBothResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SchemaObject
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case rsp.StatusCode == 200:
		// Content-type (application/octet-stream) unsupported

	}

	return response, nil
}

// ParsePostJsonResponse parses an HTTP response from a PostJsonWithResponse call
func ParsePostJsonResponse(rsp *http.Response) (*PostJsonResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &PostJsonResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetJsonResponse parses an HTTP response from a GetJsonWithResponse call
func ParseGetJsonResponse(rsp *http.Response) (*GetJsonResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetJsonResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SchemaObject
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostOtherResponse parses an HTTP response from a PostOtherWithResponse call
func ParsePostOtherResponse(rsp *http.Response) (*PostOtherResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &PostOtherResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetOtherResponse parses an HTTP response from a GetOtherWithResponse call
func ParseGetOtherResponse(rsp *http.Response) (*GetOtherResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetOtherResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetJsonWithTrailingSlashResponse parses an HTTP response from a GetJsonWithTrailingSlashWithResponse call
func ParseGetJsonWithTrailingSlashResponse(rsp *http.Response) (*GetJsonWithTrailingSlashResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetJsonWithTrailingSlashResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest SchemaObject
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /with_both_bodies)
	PostBoth(ctx echo.Context) error

	// (GET /with_both_responses)
	GetBoth(ctx echo.Context) error

	// (POST /with_json_body)
	PostJson(ctx echo.Context) error

	// (GET /with_json_response)
	GetJson(ctx echo.Context) error

	// (POST /with_other_body)
	PostOther(ctx echo.Context) error

	// (GET /with_other_response)
	GetOther(ctx echo.Context) error

	// (GET /with_trailing_slash/)
	GetJsonWithTrailingSlash(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostBoth converts echo context to params.
func (w *ServerInterfaceWrapper) PostBoth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostBoth(ctx)
	return err
}

// GetBoth converts echo context to params.
func (w *ServerInterfaceWrapper) GetBoth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetBoth(ctx)
	return err
}

// PostJson converts echo context to params.
func (w *ServerInterfaceWrapper) PostJson(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostJson(ctx)
	return err
}

// GetJson converts echo context to params.
func (w *ServerInterfaceWrapper) GetJson(ctx echo.Context) error {
	var err error

	ctx.Set("OpenId.Scopes", []string{"json.read", "json.admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetJson(ctx)
	return err
}

// PostOther converts echo context to params.
func (w *ServerInterfaceWrapper) PostOther(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostOther(ctx)
	return err
}

// GetOther converts echo context to params.
func (w *ServerInterfaceWrapper) GetOther(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetOther(ctx)
	return err
}

// GetJsonWithTrailingSlash converts echo context to params.
func (w *ServerInterfaceWrapper) GetJsonWithTrailingSlash(ctx echo.Context) error {
	var err error

	ctx.Set("OpenId.Scopes", []string{"json.read", "json.admin"})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetJsonWithTrailingSlash(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST("/with_both_bodies", wrapper.PostBoth)
	router.GET("/with_both_responses", wrapper.GetBoth)
	router.POST("/with_json_body", wrapper.PostJson)
	router.GET("/with_json_response", wrapper.GetJson)
	router.POST("/with_other_body", wrapper.PostOther)
	router.GET("/with_other_response", wrapper.GetOther)
	router.GET("/with_trailing_slash/", wrapper.GetJsonWithTrailingSlash)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9SUT2vbTBDGv4qY9z2qltPedGwPJYU2hRh6cE1Yr8bWBGl3OzNOEEbfvczaqW0aUheK",
	"IRcz6/mjZ57fSlvwsU8xYFCBegviW+xdDm9zeLO8R692ThwTshLm7IpY9Ivr0Q46JIQaRJnCGsYSOHbP",
	"JSyDPzbE2EA931WVR6MWo5VQWEVrblA8U1KKAWqYtSSFoqgUjy1qi1xoi8WHjjBoQVKsMSA7xabwkRm9",
	"dsP3ACV05DFIlhOyXvh8PTORSmoqYYaixS3yAzKU8IAsuydeTaaTqRXGhMElghreTaaTKyghOW2zDdUj",
	"aXu3jPmn2XuTomTHzC9n8q8bqOFrFH0ftYWdCWinZrA6H4NiyC0upY58bqruxWQ8MbHof8YV1PBfdYBW",
	"7YlVJ7jMxuNR0SvqG1FG15+OXEXunUINSwqOByh/Y3YCTXmD+Q9JMYitGzZdZzVHThxlt7DGZ7z4iAcr",
	"jmrfTqevwI9xPOxrmoz88DL3Tyb9Itz/ilZW/5R9CdYv/ZeBldcQ9BsmHaCeb+EmYRYzBxs8YXQNlLvY",
	"NT0FWIyLw17RPg9nYLmxurO5XOwl2sk/h8thgbPB/LOrr+yoo7C+k85JW/3p+nwjbWf7llvreAX3aRx/",
	"BgAA//+NBURdIAcAAA==",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
