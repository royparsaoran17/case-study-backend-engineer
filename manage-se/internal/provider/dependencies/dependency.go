package dependencies

import "net/http"

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type Dependency struct {
	HttpClient HttpClient
}
