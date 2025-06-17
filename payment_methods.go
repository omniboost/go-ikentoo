package ikentoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/omniboost/go-ikentoo/utils"
)

func (c *Client) NewPaymentMethodsGetRequest() PaymentMethodsGetRequest {
	r := PaymentMethodsGetRequest{
		client:  c,
		method:  http.MethodGet,
		headers: http.Header{},
	}

	r.queryParams = r.NewQueryParams()
	r.pathParams = r.NewPathParams()
	r.requestBody = r.NewRequestBody()
	return r
}

type PaymentMethodsGetRequest struct {
	client      *Client
	queryParams *PaymentMethodsGetRequestQueryParams
	pathParams  *PaymentMethodsGetRequestPathParams
	method      string
	headers     http.Header
	requestBody PaymentMethodsGetRequestBody
}

func (r PaymentMethodsGetRequest) NewQueryParams() *PaymentMethodsGetRequestQueryParams {
	return &PaymentMethodsGetRequestQueryParams{}
}

type PaymentMethodsGetRequestQueryParams struct {
	Date             Date     `schema:"date,omitempty"`
	IncludeConsumers bool     `schema:"includeConsumers,omitempty"`
	Include          Includes `schema:"include,omitempty"`
}

func (p PaymentMethodsGetRequestQueryParams) ToURLValues() (url.Values, error) {
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

func (r *PaymentMethodsGetRequest) QueryParams() *PaymentMethodsGetRequestQueryParams {
	return r.queryParams
}

func (r PaymentMethodsGetRequest) NewPathParams() *PaymentMethodsGetRequestPathParams {
	return &PaymentMethodsGetRequestPathParams{}
}

type PaymentMethodsGetRequestPathParams struct {
	BusinessID int
}

func (p *PaymentMethodsGetRequestPathParams) Params() map[string]string {
	return map[string]string{
		"business_id": strconv.Itoa(p.BusinessID),
	}
}

func (r *PaymentMethodsGetRequest) PathParams() *PaymentMethodsGetRequestPathParams {
	return r.pathParams
}

func (r *PaymentMethodsGetRequest) PathParamsInterface() PathParams {
	return r.pathParams
}

func (r *PaymentMethodsGetRequest) SetMethod(method string) {
	r.method = method
}

func (r *PaymentMethodsGetRequest) Method() string {
	return r.method
}

func (r PaymentMethodsGetRequest) NewRequestBody() PaymentMethodsGetRequestBody {
	return PaymentMethodsGetRequestBody{}
}

type PaymentMethodsGetRequestBody struct {
}

func (r *PaymentMethodsGetRequest) RequestBody() *PaymentMethodsGetRequestBody {
	return nil
}

func (r *PaymentMethodsGetRequest) RequestBodyInterface() interface{} {
	return nil
}

func (r *PaymentMethodsGetRequest) SetRequestBody(body PaymentMethodsGetRequestBody) {
	r.requestBody = body
}

func (r *PaymentMethodsGetRequest) NewResponseBody() *PaymentMethodsGetResponseBody {
	return &PaymentMethodsGetResponseBody{}
}

type PaymentMethodsGetResponseBody struct {
	Embedded struct {
		PaymentMethodList PaymentMethodList `json:"paymentMethodList"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
}

func (r *PaymentMethodsGetRequest) URL() *url.URL {
	u := r.client.GetEndpointURL("/finance/{{.business_id}}/paymentMethods", r.PathParams())
	return &u
}

func (r *PaymentMethodsGetRequest) Do(ctx context.Context) (PaymentMethodsGetResponseBody, error) {
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
