package ikentoo

import (
	"context"
	"net/http"
	"net/url"

	"github.com/omniboost/go-ikentoo/utils"
)

func (c *Client) NewBusinessesGetRequest() BusinessesGetRequest {
	r := BusinessesGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type BusinessesGetRequest struct {
	client      *Client
	queryParams *BusinessesGetRequestQueryParams
	pathParams  *BusinessesGetRequestPathParams
	method      string
	headers     http.Header
	requestBody BusinessesGetRequestBody
}

func (r BusinessesGetRequest) NewQueryParams() *BusinessesGetRequestQueryParams {
	return &BusinessesGetRequestQueryParams{}
}

type BusinessesGetRequestQueryParams struct{}

func (p BusinessesGetRequestQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	encoder.RegisterEncoder(Date{}, utils.EncodeSchemaMarshaler)
	encoder.RegisterEncoder(DateTime{}, utils.EncodeSchemaMarshaler)
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *BusinessesGetRequest) QueryParams() *BusinessesGetRequestQueryParams {
	return r.queryParams
}

func (r BusinessesGetRequest) NewPathParams() *BusinessesGetRequestPathParams {
	return &BusinessesGetRequestPathParams{}
}

type BusinessesGetRequestPathParams struct {
}

func (p *BusinessesGetRequestPathParams) Params() map[string]string {
	return map[string]string{}
}

func (r *BusinessesGetRequest) PathParams() *BusinessesGetRequestPathParams {
	return r.pathParams
}

func (r *BusinessesGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *BusinessesGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *BusinessesGetRequest) Method() string {
	return r.method
}

func (r BusinessesGetRequest) NewRequestBody() BusinessesGetRequestBody {
	return BusinessesGetRequestBody{}
}

type BusinessesGetRequestBody struct {
}

func (r *BusinessesGetRequest) RequestBody() *BusinessesGetRequestBody {
	return nil
}

func (r *BusinessesGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *BusinessesGetRequest) SetRequestBody(body BusinessesGetRequestBody) {
	r.requestBody = body
}

func (r *BusinessesGetRequest) NewResponseBody() *BusinessesGetResponseBody {
	return &BusinessesGetResponseBody{}
}

type BusinessesGetResponseBody struct {
	Embedded struct {
		BusinessList BusinessList `json:"businessList"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
}

func (r *BusinessesGetRequest) URL() *url.URL {
	u := r.client.GetEndpointURL("/data/businesses", r.PathParams())
	return &u
}

func (r *BusinessesGetRequest) Do(ctx context.Context) (BusinessesGetResponseBody, error) {
	// Create http request
	req, err := r.client.NewRequest(ctx, r)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	// Process query parameters
	err = utils.AddQueryParamsToRequest(r.QueryParams(), req, false)
	if err != nil {
		return *r.NewResponseBody(), err
	}

	responseBody := r.NewResponseBody()
	_, err = r.client.Do(req, responseBody)
	return *responseBody, err
}
