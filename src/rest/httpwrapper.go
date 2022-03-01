package rest

import "net/http"

// Interface needed to mock http Do method
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClientMock struct {
	Req *http.Request
	Res *http.Response
	Err error
}

func (c HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	return c.Res, c.Err
}
