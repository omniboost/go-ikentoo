package ikentoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/omniboost/go-ikentoo/utils"
)

func (c *Client) NewAccountingGroupsGetRequest() AccountingGroupsGetRequest {
	r := AccountingGroupsGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type AccountingGroupsGetRequest struct {
	client      *Client
	queryParams *AccountingGroupsGetRequestQueryParams
	pathParams  *AccountingGroupsGetRequestPathParams
	method      string
	headers     http.Header
	requestBody AccountingGroupsGetRequestBody
}

func (r AccountingGroupsGetRequest) NewQueryParams() *AccountingGroupsGetRequestQueryParams {
	return &AccountingGroupsGetRequestQueryParams{}
}

type AccountingGroupsGetRequestQueryParams struct {
	Include       Includes `schema:"include,omitempty"`
	PageSize      int      `schema:"pageSize,omitempty"`
	NextPageToken string   `schema:"nextPageToken,omitempty"`
}

func (p AccountingGroupsGetRequestQueryParams) ToURLValues() (url.Values, error) {
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

func (r *AccountingGroupsGetRequest) QueryParams() *AccountingGroupsGetRequestQueryParams {
	return r.queryParams
}

func (r AccountingGroupsGetRequest) NewPathParams() *AccountingGroupsGetRequestPathParams {
	return &AccountingGroupsGetRequestPathParams{}
}

type AccountingGroupsGetRequestPathParams struct {
	BusinessID int
	From       time.Time
	To         time.Time
}

func (p *AccountingGroupsGetRequestPathParams) Params() map[string]string {
	return map[string]string{
		"business_id": strconv.Itoa(p.BusinessID),
		"from":        p.From.Format("2006-01-02T15:04:05+07:00"),
		"to":          p.To.Format("2006-01-02T15:04:05+07:00"),
	}
}

func (r *AccountingGroupsGetRequest) PathParams() *AccountingGroupsGetRequestPathParams {
	return r.pathParams
}

func (r *AccountingGroupsGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *AccountingGroupsGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *AccountingGroupsGetRequest) Method() string {
	return r.method
}

func (r AccountingGroupsGetRequest) NewRequestBody() AccountingGroupsGetRequestBody {
	return AccountingGroupsGetRequestBody{}
}

type AccountingGroupsGetRequestBody struct {
}

func (r *AccountingGroupsGetRequest) RequestBody() *AccountingGroupsGetRequestBody {
	return nil
}

func (r *AccountingGroupsGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *AccountingGroupsGetRequest) SetRequestBody(body AccountingGroupsGetRequestBody) {
	r.requestBody = body
}

func (r *AccountingGroupsGetRequest) NewResponseBody() *AccountingGroupsGetResponseBody {
	return &AccountingGroupsGetResponseBody{}
}

type AccountingGroupsGetResponseBody struct {
	Embedded struct {
		AccountingGroupList AccountingGroupList `json:"accountingGroupList"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
}

func (r *AccountingGroupsGetRequest) URL() *url.URL {
	u := r.client.GetEndpointURL("/finance/{{.business_id}}/accountingGroups", r.PathParams())
	return &u
}

func (r *AccountingGroupsGetRequest) Do(ctx context.Context) (AccountingGroupsGetResponseBody, error) {
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
