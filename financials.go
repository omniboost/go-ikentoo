package ikentoo

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/omniboost/go-ikentoo/utils"
)

func (c *Client) NewFinancialsGetRequest() FinancialsGetRequest {
	r := FinancialsGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type FinancialsGetRequest struct {
	client      *Client
	queryParams *FinancialsGetRequestQueryParams
	pathParams  *FinancialsGetRequestPathParams
	method      string
	headers     http.Header
	requestBody FinancialsGetRequestBody
}

func (r FinancialsGetRequest) NewQueryParams() *FinancialsGetRequestQueryParams {
	return &FinancialsGetRequestQueryParams{}
}

type FinancialsGetRequestQueryParams struct {
	Include       Includes `schema:"include,omitempty"`
	PageSize      int      `schema:"pageSize,omitempty"`
	NextPageToken string   `schema:"nextPageToken,omitempty"`
}

func (p FinancialsGetRequestQueryParams) ToURLValues() (url.Values, error) {
	encoder := utils.NewSchemaEncoder()
	encoder.RegisterEncoder(Date{}, utils.EncodeSchemaMarshaler)
	encoder.RegisterEncoder(DateTime{}, utils.EncodeSchemaMarshaler)
	encoder.RegisterEncoder(Includes{}, utils.EncodeSchemaMarshaler)
	params := url.Values{}

	err := encoder.Encode(p, params)
	if err != nil {
		return params, err
	}

	return params, nil
}

func (r *FinancialsGetRequest) QueryParams() *FinancialsGetRequestQueryParams {
	return r.queryParams
}

func (r FinancialsGetRequest) NewPathParams() *FinancialsGetRequestPathParams {
	return &FinancialsGetRequestPathParams{}
}

type FinancialsGetRequestPathParams struct {
	BusinessID int
	From       time.Time
	To         time.Time
}

func (p *FinancialsGetRequestPathParams) Params() map[string]string {
	return map[string]string{
		"business_id": strconv.Itoa(p.BusinessID),
		"from":        p.From.Format("2006-01-02T15:04:05+07:00"),
		"to":          p.To.Format("2006-01-02T15:04:05+07:00"),
	}
}

func (r *FinancialsGetRequest) PathParams() *FinancialsGetRequestPathParams {
	return r.pathParams
}

func (r *FinancialsGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *FinancialsGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *FinancialsGetRequest) Method() string {
	return r.method
}

func (r FinancialsGetRequest) NewRequestBody() FinancialsGetRequestBody {
	return FinancialsGetRequestBody{}
}

type FinancialsGetRequestBody struct {
}

func (r *FinancialsGetRequest) RequestBody() *FinancialsGetRequestBody {
	return nil
}

func (r *FinancialsGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *FinancialsGetRequest) SetRequestBody(body FinancialsGetRequestBody) {
	r.requestBody = body
}

func (r *FinancialsGetRequest) NewResponseBody() *FinancialsGetResponseBody {
	return &FinancialsGetResponseBody{}
}

type FinancialsGetResponseBody struct {
	BusinessName       string `json:"businessName"`
	BusinessLocationID int64  `json:"businessLocationId"`
	Sales              Sales  `json:"sales"`
	DataComplete       bool   `json:"dataComplete"`
	NextPageToken      string `json:"nextPageToken"`
	Links              Links  `json:"_links"`
}

func (r *FinancialsGetRequest) URL() *url.URL {
	u := r.client.GetEndpointURL("/finance/{{.business_id}}/financials/{{.from}}/{{.to}}", r.PathParams())
	return &u
}

func (r *FinancialsGetRequest) Do() (FinancialsGetResponseBody, error) {
	// Create http request
	req, err := r.client.NewRequest(nil, r)
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

func (r *FinancialsGetRequest) All() (FinancialsGetResponseBody, error) {
	resp, err := r.Do()
	if err != nil {
		return resp, err
	}

	concat := resp

	for resp.Links.NextPage.Href != "" {
		r.QueryParams().NextPageToken = resp.NextPageToken
		resp, err = r.Do()
		if err != nil {
			return resp, err
		}

		concat.Sales = append(concat.Sales, resp.Sales...)
	}

	return concat, nil
}
