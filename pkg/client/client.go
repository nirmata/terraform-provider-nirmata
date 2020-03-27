package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
)

// RESTRequest is used to pass HTTP Request parameters
type RESTRequest struct {
	Service     Service
	Method      string
	Headers     map[string]string
	Path        string
	Data        []byte
	ContentType string
	QueryParams map[string]string
}

// Client provides access to the Nirmata REST API
type Client interface {

	// Get retrieves all data for a single Object
	Get(id ID, opts *GetOptions) (map[string]interface{}, error)

	// GetURL retrieves data from a URL path
	GetURL(service Service, urlPath string) ([]byte, int, error)

	// GetURLWithID is used to call a custom REST endpoint
	GetURLWithID(id ID, urlPath string) ([]byte, int, error)

	// GetRelationID retrieves a single relation ID of an object, by name
	GetRelationID(id ID, path string) (ID, error)

	// GetRelation retrieves a single relation of an object, by name
	GetRelation(id ID, path string) (map[string]interface{}, error)

	// GetDescendents retrieves all descendents that match a relation name or a model type
	GetDescendants(id ID, path string, opts *GetOptions) ([]map[string]interface{}, error)

	// GetAll retrieves a list of objects. The service and modelIndex are required.
	// Additional options, like which fields to retrieve and a query to filter
	// results, can be passed using opts.
	GetCollection(service Service, modelIndex string, opts *GetOptions) ([]map[string]interface{}, error)

	// QueryByName is a convinience method that queries a service collection to find
	// an object by its 'name' attribute. If a matching object is found, its ID is
	// returned.
	QueryByName(service Service, modelIndex, name string) (ID, error)

	// WaitForState queries a state and returns when it matches the specified value
	// or maxTime is reached
	WaitForState(id ID, fieldIndex string, value interface{}, maxTime time.Duration, msg string) error

	// Post is used to create a new model object.
	Post(rr *RESTRequest) (map[string]interface{}, error)

	// Post is used to create a new model object. The contentType is assumed to be JSON. The queryParams is optional.
	PostFromJSON(service Service, path string, jsonMap map[string]interface{}, queryParams map[string]string) (map[string]interface{}, error)

	// Post is used to create resources or call a custom REST endpoint. The contentType is assumed to be YAML.
	// The path and queryParams are optional
	PostWithID(id ID, path string, data []byte, queryParams map[string]string) (map[string]interface{}, error)

	// Put is used to modify a model object.
	Put(rr *RESTRequest) (map[string]interface{}, error)

	// Delete is used to delete a resource
	Delete(id ID, params map[string]string) error

	// DeleteURL is used to delete a resource identified by the URL
	DeleteURL(service Service, url string) error

	// Options is used to execute an HTTP OPTIONS request
	Options(service Service, url string) (map[string]interface{}, error)
}

// GetOptions contains optional paramameters used when retrieving objects
type GetOptions struct {
	Fields []string
	Filter Query
	Mode   OutputMode
}

// NewGetOptions build a new GetOptions instance
func NewGetOptions(fields []string, query Query) *GetOptions {
	return &GetOptions{fields, query, OutputModeNone}
}

// NewGetModelID builds a new GetOptions instance with Model ID attribute fields included
func NewGetModelID(fields []string, query Query) *GetOptions {
	allFields := []string{"id", "modelIndex", "service", "name"}
	if fields != nil {
		allFields = append(allFields, fields...)
	}

	return &GetOptions{allFields, query, OutputModeNone}
}

// NewClient creates a new Client
func NewClient(address string, token string, httpClient *http.Client, insecure bool) Client {
	if insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &client{address: address, token: token, httpClient: httpClient}
}

type client struct {
	address    string
	token      string
	httpClient *http.Client
}

func (c *client) GetURLWithID(id ID, urlPath string) ([]byte, int, error) {

	rawURL := removeSlash(c.address) + "/" +
		id.Service().Name() + "/api/" +
		id.ModelIndex() + "/" +
		id.UUID() + "/" +
		escapeQuery(urlPath)

	return c.get(rawURL)
}

func (c *client) GetURL(service Service, urlPath string) ([]byte, int, error) {
	p := strings.Join([]string{removeSlash(c.address), service.Name(), "api", escapeQuery(urlPath)}, "/")
	return c.get(p)
}

func removeSlash(path string) string {
	return strings.TrimRight(path, "/")
}

func (c *client) get(rawURL string) ([]byte, int, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))
	glog.V(3).Infof("HTTP %s request %s", req.Method, req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		glog.V(1).Infof("HTTP %d '%s': %s", resp.StatusCode, resp.Status, string(b))
		return b, resp.StatusCode, fmt.Errorf("HTTP %s", resp.Status)
	}

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	return b, resp.StatusCode, nil
}

func escapeQuery(path string) string {
	parts := strings.Split(path, "?")
	if len(parts) != 2 {
		return path
	}

	params, err := url.ParseQuery(parts[1])
	if err != nil {
		return path
	}

	return parts[0] + "?" + params.Encode()
}

func (c *client) Get(id ID, opts *GetOptions) (map[string]interface{}, error) {
	ubldr := NewURLBuilder(c.address).
		ToService(id.Service()).
		WithPaths(id.ModelIndex(), id.UUID())

	if opts != nil {
		ubldr.SelectFields(opts.Fields)
		ubldr.WithQuery(opts.Filter)
		ubldr.WithMode(opts.Mode)
	}

	u := ubldr.Build()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))

	glog.V(3).Infof("HTTP %s request %s", req.Method, req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ParseObject(b)
	}

	return nil, fmt.Errorf("%s: %s", resp.Status, string(b))
}

func (c *client) GetRelationID(id ID, name string) (ID, error) {
	b, err := c.getRelationData(id, name)
	if err != nil {
		return nil, err
	}

	return ParseID(b)
}

func (c *client) GetRelation(id ID, name string) (map[string]interface{}, error) {
	b, err := c.getRelationData(id, name)
	if err != nil {
		return nil, err
	}

	data, err := ParseCollection(b)
	if err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, nil
	}

	return data[0], nil
}

func (c *client) getRelationData(id ID, name string) ([]byte, error) {
	uBldr := NewURLBuilder(c.address).
		ToService(id.Service()).
		WithPaths(id.ModelIndex(), id.UUID(), name)

	u := uBldr.Build()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))

	glog.V(3).Infof("HTTP %s request %s", req.Method, req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return b, nil
	}

	return nil, fmt.Errorf("%s: %s", resp.Status, string(b))

}

func (c *client) GetDescendants(id ID, path string, opts *GetOptions) ([]map[string]interface{}, error) {

	uBldr := NewURLBuilder(c.address).
		ToService(id.Service()).
		WithPaths(id.ModelIndex(), id.UUID(), path)

	if opts != nil {
		uBldr.SelectFields(opts.Fields)
		uBldr.WithQuery(opts.Filter)
	}

	u := uBldr.Build()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))

	glog.V(3).Infof("HTTP %s request %s", req.Method, req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ParseCollection(b)
	}

	return nil, fmt.Errorf("%s: %s", resp.Status, string(b))
}

func (c *client) Delete(id ID, params map[string]string) error {
	u := NewURLBuilder(c.address).
		ToService(id.Service()).
		WithPaths(id.ModelIndex(), id.UUID()).
		WithParameters(params).
		Build()

	return c.delete(u)
}

func (c *client) DeleteURL(service Service, path string) error {
	u := NewURLBuilder(c.address).
		ToService(service).
		WithPaths(path).
		Build()

	return c.delete(u)
}

func (c *client) delete(u string) error {
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))

	glog.V(3).Infof("HTTP %s request %s", req.Method, req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("%s: %s", resp.Status, string(b))
}

func (c *client) GetCollection(service Service, modelIndex string, opts *GetOptions) ([]map[string]interface{}, error) {

	ubldr := NewURLBuilder(c.address).
		ToService(service).
		WithPath(modelIndex)

	if opts != nil {
		ubldr.SelectFields(opts.Fields)
		ubldr.WithQuery(opts.Filter)
	}

	u := ubldr.Build()
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))

	glog.V(3).Infof("HTTP %s request %s", req.Method, req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ParseCollection(b)
	}

	return nil, fmt.Errorf("%s: %s", resp.Status, string(b))
}

func (c *client) Post(rr *RESTRequest) (map[string]interface{}, error) {
	req, err := c.buildRequest("POST", rr)
	if err != nil {
		return nil, err
	}

	return c.send(req)
}

func (c *client) Put(rr *RESTRequest) (map[string]interface{}, error) {
	req, err := c.buildRequest("PUT", rr)
	if err != nil {
		return nil, err
	}

	return c.send(req)
}

func (c *client) Options(service Service, url string) (map[string]interface{}, error) {
	rr := &RESTRequest{
		Service: service,
		Path:    url,
	}

	req, err := c.buildRequest("OPTIONS", rr)
	if err != nil {
		return nil, err
	}

	return c.send(req)
}

func (c *client) buildRequest(method string, rr *RESTRequest) (*http.Request, error) {
	u := NewURLBuilder(c.address).
		ToService(rr.Service).
		WithPath(rr.Path).
		Build()

	b := rr.Data
	if b == nil {
		b = []byte{}
	}

	req, err := http.NewRequest(method, u, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	// set default headers before user supplied headers
	req.Header.Add("Authorization", fmt.Sprintf("NIRMATA-API %s", c.token))
	if rr.ContentType != "" {
		req.Header.Set("Content-Type", rr.ContentType)
	}

	// Overrride default headers with user supplied headers
	if rr.Headers != nil {
		for k, v := range rr.Headers {
			req.Header.Set(k, v)
		}
	}

	if rr.QueryParams != nil {
		q := req.URL.Query()
		for k, v := range rr.QueryParams {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (c *client) send(request *http.Request) (map[string]interface{}, error) {
	glog.V(3).Infof("HTTP %s request %s", request.Method, request.URL.String())
	glog.V(10).Infof("HTTP request details:\n %s", dumpRequest(request))

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	glog.V(3).Infof("HTTP response %s - body[%d bytes]", resp.Status, len(b))
	glog.V(10).Infof("HTTP response body: %s", string(b))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ParseObject(b)
	}

	return nil, fmt.Errorf("%s: %s", resp.Status, string(b))
}

func dumpRequest(request *http.Request) string {
	requestDump, _ := httputil.DumpRequest(request, true)
	return string(requestDump)
}

func (c *client) PostFromJSON(service Service, path string, jsonMap map[string]interface{}, queryParams map[string]string) (map[string]interface{}, error) {
	b, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}

	req := &RESTRequest{
		Service:     service,
		Path:        path,
		ContentType: "application/json",
		QueryParams: queryParams,
		Data:        b,
	}

	return c.Post(req)
}

func (c *client) PostWithID(id ID, path string, data []byte, queryParams map[string]string) (map[string]interface{}, error) {
	service := id.Service()
	requestPath := id.ModelIndex() + "/" + id.UUID()
	if path != "" {
		requestPath = requestPath + "/" + path
	}

	req := &RESTRequest{
		Service:     service,
		Path:        requestPath,
		ContentType: "application/yaml",
		QueryParams: queryParams,
		Data:        data,
	}

	return c.Post(req)
}

func (c *client) QueryByName(service Service, modelIndex, name string) (ID, error) {
	opts := &GetOptions{}
	opts.Filter = NewQuery().FieldEqualsValue("name", name)
	opts.Fields = []string{"id", "name", "modelIndex", "service"}

	objs, err := c.GetCollection(service, modelIndex, opts)
	if err != nil {
		return nil, err
	}

	if len(objs) == 0 {
		return nil, fmt.Errorf("Failed to find %s with name %s in service %s",
			modelIndex, name, service.Name())
	}

	if len(objs) > 1 {
		return nil, fmt.Errorf("Multiple %s instances with name %s in service %s",
			modelIndex, name, service.Name())
	}

	obj, err := NewObject(objs[0])
	if err != nil {
		return nil, err
	}

	return obj.ID(), nil
}

func (c *client) WaitForState(id ID, fieldIndex string, value interface{}, maxTime time.Duration, msg string) error {

	timer := time.NewTimer(maxTime)
	defer timer.Stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		fmt.Print(msg)
		select {

		case <-timer.C:
			return fmt.Errorf("Timed out on %s = %v", fieldIndex, value)

		case <-ticker.C:
			match, err := c.checkState(id, fieldIndex, value)
			if err != nil {
				return err
			}

			if match {
				return nil
			}
		}
	}
}

func (c *client) checkState(id ID, fieldIndex string, value interface{}) (bool, error) {
	data, err := c.Get(id, NewGetModelID([]string{fieldIndex}, nil))
	if err != nil {
		return false, err
	}

	o, err := NewObject(data)
	if err != nil {
		return false, err
	}

	rval := o.Data()[fieldIndex]
	glog.V(1).Infof("%s %s = %v\n", id.ModelIndex(), fieldIndex, rval)
	if rval != nil && rval == value {
		return true, nil
	}

	return false, nil
}
