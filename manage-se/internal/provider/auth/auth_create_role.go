package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"manage-se/internal/presentations"
	"manage-se/internal/provider/providererrors"
	"manage-se/pkg/httpx"
	"net/http"
)

func (c *client) CreateRole(ctx context.Context, input presentations.RoleCreate) (*Role, error) {
	urlEndpoint := c.endpoint("/internal/v1/roles")

	var request bytes.Buffer
	err := json.NewEncoder(&request).Encode(input)
	if err != nil {
		return nil, errors.Wrap(err, "new encoder encode")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlEndpoint, &request)
	if err != nil {
		return nil, errors.Wrap(err, "new request with context")
	}

	req.Header.Set(httpx.ContentType, httpx.MediaTypeJSON)

	res, err := c.dep.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	rawBody, _ := io.ReadAll(res.Body)
	res.Body.Close() // must close
	res.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	switch res.StatusCode {
	case http.StatusOK:
		body := struct {
			Data Role `json:"data"`
		}{}

		err = json.Unmarshal(rawBody, &body)
		if err != nil {
			return nil, providererrors.NewErrRequestWithResponse(req, res)
		}

		return &body.Data, nil

	default:
		bodyErr := providererrors.Error{}
		err := json.Unmarshal(rawBody, &bodyErr)
		if err != nil {
			return nil, providererrors.NewErrRequestWithResponse(req, res)
		}

		bodyErr.Code = res.StatusCode
		return nil, bodyErr

	}
}
