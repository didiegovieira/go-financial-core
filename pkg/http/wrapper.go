package http

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Wrapper interface {
	Get(ctx context.Context, url string, request Request) (*Response, error)
	Post(ctx context.Context, url string, request Request) (*Response, error)
	Put(ctx context.Context, url string, request Request) (*Response, error)
	Patch(ctx context.Context, url string, request Request) (*Response, error)
	Delete(ctx context.Context, url string, request Request) (*Response, error)
	Request(ctx context.Context, method string, url string, request Request) (*Response, error)
}

type Request struct {
	Headers map[string]string
	Body    []byte
}

type Response struct {
	Headers    map[string]string
	Body       []byte
	StatusCode int
}

type WrapperImpl struct {
	Client *http.Client
}

func NewWrapper() *WrapperImpl {
	return &WrapperImpl{
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			},
		},
	}
}

func (d *WrapperImpl) Get(ctx context.Context, url string, request Request) (*Response, error) {
	return d.Request(ctx, http.MethodGet, url, request)
}

func (d *WrapperImpl) Post(ctx context.Context, url string, request Request) (*Response, error) {
	return d.Request(ctx, http.MethodPost, url, request)
}

func (d *WrapperImpl) Put(ctx context.Context, url string, request Request) (*Response, error) {
	return d.Request(ctx, http.MethodPut, url, request)
}

func (d *WrapperImpl) Patch(ctx context.Context, url string, request Request) (*Response, error) {
	return d.Request(ctx, http.MethodPatch, url, request)
}

func (d *WrapperImpl) Delete(ctx context.Context, url string, request Request) (*Response, error) {
	return d.Request(ctx, http.MethodDelete, url, request)
}

func (d *WrapperImpl) Request(ctx context.Context, method string, url string, request Request) (*Response, error) {
	var err error
	var req *http.Request
	var payload io.Reader = nil

	if request.Body != nil {
		payload = bytes.NewBuffer(request.Body)
	}

	if ctx == context.TODO() {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, payload)
	}

	if err != nil {
		return nil, err
	}

	if len(request.Headers) > 0 {
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	for key, value := range resp.Header {
		headers[key] = value[0]
	}

	response := Response{Headers: headers, Body: body, StatusCode: resp.StatusCode}

	return &response, nil
}

func (r *Response) IsSuccessful() bool {
	return r.StatusCode == http.StatusOK || r.StatusCode == http.StatusCreated
}

func (r *Response) Error() error {
	if r.IsSuccessful() {
		return nil
	}
	return errors.New("Http error " + strconv.Itoa(r.StatusCode) + ": " + string(r.Body))
}
