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

type DoFn func(r *http.Request) (*http.Response, error)

// RoundTripMiddleware lets you define functions that can intercept and manipulate
// the round trip of a single HTTP transaction
type RoundTripMiddleware func(next DoFn) DoFn

// DoFn returns the result of applying the middleware to the provided DoFn
func (rtm RoundTripMiddleware) DoFn(doFn DoFn) DoFn {
	return rtm(doFn)
}

func joinMiddleware(mw ...RoundTripMiddleware) RoundTripMiddleware {
	if len(mw) < 1 {
		return func(doFn DoFn) DoFn {
			return doFn
		}
	}
	middleware := mw[len(mw)-1]
	for i := len(mw) - 2; i >= 0; i-- {
		middleware = middleware.Wrap(mw[i])
	}
	return middleware
}

func (mw RoundTripMiddleware) Wrap(wrapMw RoundTripMiddleware) RoundTripMiddleware {
	return func(doFn DoFn) DoFn {
		return wrapMw(mw(doFn))
	}
}

// RoundTripMiddlewares allows configuring a RoundTripMiddleware for individual endpoints
type RoundTripMiddlewares struct {
	PostBoth                 RoundTripMiddleware
	GetBoth                  RoundTripMiddleware
	PostJson                 RoundTripMiddleware
	GetJson                  RoundTripMiddleware
	PostOther                RoundTripMiddleware
	GetOther                 RoundTripMiddleware
	GetJsonWithTrailingSlash RoundTripMiddleware
}

// operationDoFunctions lets the client store Do functions using different
// middleware for each operation
type operationDoFunctions struct {
	PostBoth                 DoFn
	GetBoth                  DoFn
	PostJson                 DoFn
	GetJson                  DoFn
	PostOther                DoFn
	GetOther                 DoFn
	GetJsonWithTrailingSlash DoFn
}

// client which conforms to the OpenAPI3 specification for this service.
type client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn

	// SharedRoundTripMiddleware lets you apply a RoundTripMiddleware on all
	// operations.
	SharedRoundTripMiddleware RoundTripMiddleware

	// RoundTripMiddlewares lets you apply a RoundTripMiddleware on specific
	// operations.
	RoundTripMiddlewares RoundTripMiddlewares

	// operationDoers is the set of Do functions for each operation that is created
	// for the client.
	operationDoers *operationDoFunctions
}

var _ clientInterface = &client{}

// clientOption allows setting custom parameters during construction
type clientOption func(*client) error

// newClient Creates a new Client, with reasonable defaults
func newClient(server string, opts ...clientOption) (*client, error) {
	// create a client with sane default values
	client := client{
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

	client.operationDoers = setupOperationDoers(&client, client.RoundTripMiddlewares)

	return &client, nil
}

func setupOperationDoers(c *client, rtMiddlewares RoundTripMiddlewares) *operationDoFunctions {

	sharedMiddlewares := []RoundTripMiddleware{}

	if c.RequestEditor != nil {
		mw := newRequestEditorMiddleware(c.RequestEditor)
		sharedMiddlewares = append(sharedMiddlewares, mw)
	}

	if c.SharedRoundTripMiddleware != nil {
		sharedMiddlewares = append(sharedMiddlewares, c.SharedRoundTripMiddleware)
	}

	sharedMiddleware := joinMiddleware(sharedMiddlewares...)

	operationDoers := operationDoFunctions{}

	// PostBoth
	if rtMiddlewares.PostBoth != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.PostBoth)
		operationDoers.PostBoth = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.PostBoth = sharedMiddleware.DoFn(c.Client.Do)
	}
	// GetBoth
	if rtMiddlewares.GetBoth != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.GetBoth)
		operationDoers.GetBoth = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.GetBoth = sharedMiddleware.DoFn(c.Client.Do)
	}
	// PostJson
	if rtMiddlewares.PostJson != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.PostJson)
		operationDoers.PostJson = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.PostJson = sharedMiddleware.DoFn(c.Client.Do)
	}
	// GetJson
	if rtMiddlewares.GetJson != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.GetJson)
		operationDoers.GetJson = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.GetJson = sharedMiddleware.DoFn(c.Client.Do)
	}
	// PostOther
	if rtMiddlewares.PostOther != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.PostOther)
		operationDoers.PostOther = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.PostOther = sharedMiddleware.DoFn(c.Client.Do)
	}
	// GetOther
	if rtMiddlewares.GetOther != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.GetOther)
		operationDoers.GetOther = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.GetOther = sharedMiddleware.DoFn(c.Client.Do)
	}
	// GetJsonWithTrailingSlash
	if rtMiddlewares.GetJsonWithTrailingSlash != nil {
		mw := joinMiddleware(sharedMiddleware, rtMiddlewares.GetJsonWithTrailingSlash)
		operationDoers.GetJsonWithTrailingSlash = mw.DoFn(c.Client.Do)
	} else {
		operationDoers.GetJsonWithTrailingSlash = sharedMiddleware.DoFn(c.Client.Do)
	}

	return &operationDoers
}

func newRequestEditorMiddleware(requestEditorFn RequestEditorFn) RoundTripMiddleware {
	return func(next DoFn) DoFn {
		return func(r *http.Request) (*http.Response, error) {
			err := requestEditorFn(r.Context(), r)
			if err != nil {
				return nil, err
			}
			return next(r)
		}
	}
}

// WithSharedRoundTripMiddleware add a middleware that applies to all routes
func WithSharedRoundTripMiddleware(rtm RoundTripMiddleware) clientOption {
	return func(c *Client) error {
		c.SharedRoundTripMiddleware = rtm
		return nil
	}
}

// WithRoundTripMiddlewares Add middlewares that apply to specific routes
func WithRoundTripMiddlewares(rtMiddlewares RoundTripMiddlewares) clientOption {
	return func(c *Client) error {
		c.RoundTripMiddlewares = rtMiddlewares
		return nil
	}
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) clientOption {
	return func(c *client) error {
		c.Client = doer
		return nil
	}
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) clientOption {
	return func(c *client) error {
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
func WithRequestEditorFn(fn RequestEditorFn) clientOption {
	return func(c *client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type clientInterface interface {
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

func (c *client) PostBothWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewPostBothRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.PostBoth(req)
}

func (c *client) PostBoth(ctx context.Context, body PostBothJSONRequestBody) (*http.Response, error) {
	req, err := NewPostBothRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.PostBoth(req)
}

func (c *client) GetBoth(ctx context.Context) (*http.Response, error) {
	req, err := NewGetBothRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.GetBoth(req)
}

func (c *client) PostJsonWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewPostJsonRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.PostJson(req)
}

func (c *client) PostJson(ctx context.Context, body PostJsonJSONRequestBody) (*http.Response, error) {
	req, err := NewPostJsonRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.PostJson(req)
}

func (c *client) GetJson(ctx context.Context) (*http.Response, error) {
	req, err := NewGetJsonRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.GetJson(req)
}

func (c *client) PostOtherWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewPostOtherRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.PostOther(req)
}

func (c *client) GetOther(ctx context.Context) (*http.Response, error) {
	req, err := NewGetOtherRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.GetOther(req)
}

func (c *client) GetJsonWithTrailingSlash(ctx context.Context) (*http.Response, error) {
	req, err := NewGetJsonWithTrailingSlashRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.operationDoers.GetJsonWithTrailingSlash(req)
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
func (c *client) PostBothWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostBothResponse, error) {
	rsp, err := c.PostBothWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsePostBothResponse(rsp)
}

func (c *client) PostBothWithResponse(ctx context.Context, body PostBothJSONRequestBody) (*PostBothResponse, error) {
	rsp, err := c.PostBoth(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsePostBothResponse(rsp)
}

// GetBothWithResponse request returning *GetBothResponse
func (c *client) GetBothWithResponse(ctx context.Context) (*GetBothResponse, error) {
	rsp, err := c.GetBoth(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetBothResponse(rsp)
}

// PostJsonWithBodyWithResponse request with arbitrary body returning *PostJsonResponse
func (c *client) PostJsonWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostJsonResponse, error) {
	rsp, err := c.PostJsonWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsePostJsonResponse(rsp)
}

func (c *client) PostJsonWithResponse(ctx context.Context, body PostJsonJSONRequestBody) (*PostJsonResponse, error) {
	rsp, err := c.PostJson(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParsePostJsonResponse(rsp)
}

// GetJsonWithResponse request returning *GetJsonResponse
func (c *client) GetJsonWithResponse(ctx context.Context) (*GetJsonResponse, error) {
	rsp, err := c.GetJson(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetJsonResponse(rsp)
}

// PostOtherWithBodyWithResponse request with arbitrary body returning *PostOtherResponse
func (c *client) PostOtherWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*PostOtherResponse, error) {
	rsp, err := c.PostOtherWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParsePostOtherResponse(rsp)
}

// GetOtherWithResponse request returning *GetOtherResponse
func (c *client) GetOtherWithResponse(ctx context.Context) (*GetOtherResponse, error) {
	rsp, err := c.GetOther(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetOtherResponse(rsp)
}

// GetJsonWithTrailingSlashWithResponse request returning *GetJsonWithTrailingSlashResponse
func (c *client) GetJsonWithTrailingSlashWithResponse(ctx context.Context) (*GetJsonWithTrailingSlashResponse, error) {
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

	"H4sIAAAAAAAC/9SUz4oTQRDGX2UoPY6ZrN7mqAdZQVdMwEMMS6enku5lprutquwyhHl3qU7WTHBZI0hg",
	"L6E69Yevvl9P78DGLsWAQRjqHbB12JkcznJ4s7pDK3pOFBOSeMzZtSeWL6ZDPUifEGpgIR82MJRAsX0q",
	"oRn8ufWEDdSLfVU5GrUctMSHddTmBtmST+JjgBrmznMhyMLFg0NxSIU4LD60HoMUJjSH8LsX9w05xcDI",
	"hSEsNhiQjGBT2EiEVtr+R4ASWm8xcNYZ8iLw+Xqu6sWLyoc5shQzpHskKOEeifdSribTyVQLY8Jgkoca",
	"3k2mkysoIRlx2Z/qwYu7XcX80xxMS5GzlWqk0b2uG6jha2R5H8XB3h3UU9NrnY1BMOQWk1LrbW6q7lhl",
	"PMLS6DXhGmp4VR1pVgeU1QlH9Xc8KlpBecNCaLrTketInRGoYeWDoR7KP2Ce0BTaYv7j4DzUYdu2WjNy",
	"YpTdwQaf8OIjHq0Y1b6dTl+AH8Nw3Fc1Kfn+ee6fVPpFuP8Traz+MfscrN/6LwMrr8Fot+Slh3qxg5uE",
	"WcwCdPCE0DRQ7mPTdD7Aclge94r6bpyB5UbrzuZysY9oL/8cLscFzgbz366+kPGtD5tbbg276m/XRx/s",
	"+aFlph0v4D4Nw68AAAD//7fdvjU5BwAA",
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
