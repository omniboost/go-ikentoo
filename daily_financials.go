package ikentoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/omniboost/go-ikentoo/utils"
)

func (c *Client) NewDailyFinancialsGetRequest() DailyFinancialsGetRequest {
	r := DailyFinancialsGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type DailyFinancialsGetRequest struct {
	client      *Client
	queryParams *DailyFinancialsGetRequestQueryParams
	pathParams  *DailyFinancialsGetRequestPathParams
	method      string
	headers     http.Header
	requestBody DailyFinancialsGetRequestBody
}

func (r DailyFinancialsGetRequest) NewQueryParams() *DailyFinancialsGetRequestQueryParams {
	return &DailyFinancialsGetRequestQueryParams{}
}

type DailyFinancialsGetRequestQueryParams struct {
	Date             Date     `schema:"date,omitempty"`
	IncludeConsumers bool     `schema:"includeConsumers,omitempty"`
	Include          Includes `schema:"include,omitempty"`
}

func (p DailyFinancialsGetRequestQueryParams) ToURLValues() (url.Values, error) {
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

func (r *DailyFinancialsGetRequest) QueryParams() *DailyFinancialsGetRequestQueryParams {
	return r.queryParams
}

func (r DailyFinancialsGetRequest) NewPathParams() *DailyFinancialsGetRequestPathParams {
	return &DailyFinancialsGetRequestPathParams{}
}

type DailyFinancialsGetRequestPathParams struct {
	BusinessID int
}

func (p *DailyFinancialsGetRequestPathParams) Params() map[string]string {
	return map[string]string{
		"business_id": strconv.Itoa(p.BusinessID),
	}
}

func (r *DailyFinancialsGetRequest) PathParams() *DailyFinancialsGetRequestPathParams {
	return r.pathParams
}

func (r *DailyFinancialsGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *DailyFinancialsGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *DailyFinancialsGetRequest) Method() string {
	return r.method
}

func (r DailyFinancialsGetRequest) NewRequestBody() DailyFinancialsGetRequestBody {
	return DailyFinancialsGetRequestBody{}
}

type DailyFinancialsGetRequestBody struct {
}

func (r *DailyFinancialsGetRequest) RequestBody() *DailyFinancialsGetRequestBody {
	return nil
}

func (r *DailyFinancialsGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *DailyFinancialsGetRequest) SetRequestBody(body DailyFinancialsGetRequestBody) {
	r.requestBody = body
}

func (r *DailyFinancialsGetRequest) NewResponseBody() *DailyFinancialsGetResponseBody {
	return &DailyFinancialsGetResponseBody{}
}

type DailyFinancialsGetResponseBody struct {
	BusinessName            string    `json:"businessName"`
	NextStartOfDayAsIso8601 time.Time `json:"nextStartOfDayAsIso8601"`
	BusinessLocationID      int64     `json:"businessLocationId"`
	Sales                   Sales     `json:"sales"`
	DataComplete            bool      `json:"dataComplete"`
	Links                   Links     `json:"_links"`
}

func (r *DailyFinancialsGetRequest) URL() *url.URL {
	u := r.client.GetEndpointURL("/finance/{{.business_id}}/dailyFinancials", r.PathParams())
	return &u
}

func (r *DailyFinancialsGetRequest) Do(ctx context.Context) (DailyFinancialsGetResponseBody, error) {
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
