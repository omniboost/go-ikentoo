package ikentoo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
	"text/template"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	libraryVersion = "0.0.1"
	userAgent      = "go-ikentoo/" + libraryVersion
	mediaType      = "application/json"
	charset        = "utf-8"
)

var (
	BaseURL = url.URL{
		Scheme: "https",
		Host:   "api.ikentoo.com",
		Path:   "/f",
	}
)

// NewClient returns a new Exact Globe Client client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client := &Client{}

	client.SetHTTPClient(httpClient)
	client.SetBaseURL(BaseURL)
	client.SetDebug(false)
	client.SetUserAgent(userAgent)
	client.SetMediaType(mediaType)
	client.SetCharset(charset)

	return client
}

// Client manages communication with Exact Globe Client
type Client struct {
	// HTTP client used to communicate with the Client.
	http *http.Client

	debug   bool
	baseURL url.URL

	// credentials

	// User agent for client
	userAgent string

	mediaType             string
	charset               string
	disallowUnknownFields bool

	// Optional function called after every successful request made to the DO Clients
	beforeRequestDo    BeforeRequestDoCallback
	onRequestCompleted RequestCompletionCallback
}

type BeforeRequestDoCallback func(*http.Client, *http.Request, interface{})

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

func (c *Client) SetHTTPClient(client *http.Client) {
	c.http = client
}

func (c Client) Debug() bool {
	return c.debug
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c Client) BaseURL() url.URL {
	return c.baseURL
}

func (c *Client) SetBaseURL(baseURL url.URL) {
	c.baseURL = baseURL
}

func (c *Client) SetMediaType(mediaType string) {
	c.mediaType = mediaType
}

func (c Client) MediaType() string {
	return mediaType
}

func (c *Client) SetCharset(charset string) {
	c.charset = charset
}

func (c Client) Charset() string {
	return charset
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c Client) UserAgent() string {
	return userAgent
}

func (c *Client) SetDisallowUnknownFields(disallowUnknownFields bool) {
	c.disallowUnknownFields = disallowUnknownFields
}

func (c *Client) SetBeforeRequestDo(fun BeforeRequestDoCallback) {
	c.beforeRequestDo = fun
}

func (c *Client) GetEndpointURL(p string, pathParams PathParams) url.URL {
	clientURL := c.BaseURL()

	parsed, err := url.Parse(p)
	if err != nil {
		log.Fatal(err)
	}
	q := clientURL.Query()
	for k, vv := range parsed.Query() {
		for _, v := range vv {
			q.Add(k, v)
		}
	}
	clientURL.RawQuery = q.Encode()

	clientURL.Path = path.Join(clientURL.Path, parsed.Path)

	tmpl, err := template.New("path").Parse(clientURL.Path)
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	params := pathParams.Params()
	// params["administration_id"] = c.Administration()
	err = tmpl.Execute(buf, params)
	if err != nil {
		log.Fatal(err)
	}

	clientURL.Path = buf.String()
	return clientURL
}

func (c *Client) NewRequest(ctx context.Context, req Request) (*http.Request, error) {
	// convert body struct to json
	buf := new(bytes.Buffer)
	if req.RequestBodyInterface() != nil {
		err := json.NewEncoder(buf).Encode(req.RequestBodyInterface())
		if err != nil {
			return nil, err
		}
	}

	// create new http request
	r, err := http.NewRequest(req.Method(), req.URL().String(), buf)
	if err != nil {
		return nil, err
	}

	// values := url.Values{}
	// err = utils.AddURLValuesToRequest(values, req, true)
	// if err != nil {
	// 	return nil, err
	// }

	// optionally pass along context
	if ctx != nil {
		r = r.WithContext(ctx)
	}

	// set other headers
	r.Header.Add("Content-Type", fmt.Sprintf("%s; charset=%s", c.MediaType(), c.Charset()))
	r.Header.Add("Accept", c.MediaType())
	r.Header.Add("User-Agent", c.UserAgent())

	return r, nil
}

// Do sends an Client request and returns the Client response. The Client response is json decoded and stored in the value
// pointed to by v, or returned as an error if an Client error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, body interface{}) (*http.Response, error) {
	if c.beforeRequestDo != nil {
		c.beforeRequestDo(c.http, req, body)
	}

	if c.debug {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Println(string(dump))
	}

	httpResp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, httpResp)
	}

	// close body io.Reader
	defer func() {
		if rerr := httpResp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if c.debug {
		dump, _ := httputil.DumpResponse(httpResp, true)
		log.Println(string(dump))
	}

	// check if the response isn't an error
	err = CheckResponse(httpResp)
	if err != nil {
		return httpResp, err
	}

	// check the provided interface parameter
	if httpResp == nil {
		return httpResp, nil
	}

	if body == nil {
		return httpResp, err
	}

	if httpResp.ContentLength == 0 {
		return httpResp, nil
	}

	errResp := &ErrorResponse{Response: httpResp}
	errsResp := &ErrorsResponse{Response: httpResp}
	err = c.Unmarshal(httpResp.Body, body, errResp, errsResp)
	if err != nil {
		return httpResp, err
	}

	if errResp.Error() != "" {
		return httpResp, errResp
	}

	if errsResp.Error() != "" {
		return httpResp, errsResp
	}

	return httpResp, nil
}

func (c *Client) Unmarshal(r io.Reader, vv ...interface{}) error {
	if len(vv) == 0 {
		return nil
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	errs := []error{}
	for _, v := range vv {
		r := bytes.NewReader(b)
		dec := json.NewDecoder(r)
		if c.disallowUnknownFields {
			dec.DisallowUnknownFields()
		}

		err := dec.Decode(v)
		if err != nil && err != io.EOF {
			errs = append(errs, err)
		}

	}

	if len(errs) == len(vv) {
		// Everything errored
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = fmt.Sprint(e)
		}
		return errors.New(strings.Join(msgs, ", "))
	}

	return nil
}

// CheckResponse checks the Client response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. Client error responses are expected to have either no response
// body, or a json response body that maps to ErrorResponse. Any other response
// body will be silently ignored.
func CheckResponse(r *http.Response) error {
	errorResponse := &ErrorResponse{Response: r}

	// Don't check content-lenght: a created response, for example, has no body
	// if r.Header.Get("Content-Length") == "0" {
	// 	errorResponse.Errors.Message = r.Status
	// 	return errorResponse
	// }

	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	// read data and copy it back
	data, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(data))
	if err != nil {
		return errorResponse
	}

	err = checkContentType(r)
	if err != nil {
		return errors.WithStack(err)
	}

	if r.ContentLength == 0 {
		return errors.New("response body is empty")
	}

	// convert json to struct
	if len(data) != 0 {
		err = json.Unmarshal(data, &errorResponse)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if errorResponse.Error() != "" {
		return errorResponse
	}

	return nil
}

// {
//   "timestamp": "2023-05-12T13:29:05.028+0000",
//   "status": 404,
//   "error": "Not Found",
//   "exception": "com.ikentoo.financial.service.exceptions.ResourceNotFoundException",
//   "message": "Business Location Not found for 12345",
//   "path": "/f/finance/12345/financials/2023-05-10T09:00:00+02:00/2023-05-12T09:00:00+02:00"
// }

type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Timestamp DateTime    `json:"timestamp"`
	Status    int    `json:"status"`
	Err       string `json:"error"`
	Exception string `json:"exception"`
	Message   string `json:"message"`
	Path      string `json:"path"`
}

func (r *ErrorResponse) Error() string {
	if r.Err == "" && r.Message == "" {
		return ""
	}

	return fmt.Sprintf("%s (%d): %s", r.Err, r.Status, r.Message)
}

// {
//   "errors": [
//     {
//       "code": "unauthorized",
//       "title": "Full authentication is required to access this resource"
//     }
//   ]
// }

type ErrorsResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Timestamp DateTime    `json:"timestamp"`
	Errors []struct {
		Code  string `json:"code"`
		Title string `json:"title"`
	}
}

func (r *ErrorsResponse) Error() string {
	if len(r.Errors) == 0 {
		return ""
	}

	var errs *multierror.Error
	for _, e := range r.Errors {
		m := fmt.Sprintf("%s: %s", e.Code, e.Title)
		errs = multierror.Append(errs, errors.New(m))
	}

	return errs.Error()

}
func checkContentType(response *http.Response) error {
	header := response.Header.Get("Content-Type")
	contentType := strings.Split(header, ";")[0]
	if contentType != mediaType {
		return fmt.Errorf("Expected Content-Type \"%s\", got \"%s\"", mediaType, contentType)
	}

	return nil
}
