package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"manage-se/internal/provider/providererrors"
	"net/http"
)

func (c *client) GetListTransactions(ctx context.Context) ([]Transaction, error) {
	urlEndpoint := c.endpoint("/internal/v1/transactions")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEndpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request failed")
	}

	res, err := c.dep.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("doing http request to %s", req.URL))
	}

	// Re-usable response body for logging
	rawBody, _ := io.ReadAll(res.Body)
	res.Body.Close() // must close
	res.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	switch res.StatusCode {
	case http.StatusCreated, http.StatusOK:
		body := struct {
			Data []Transaction `json:"data"`
		}{}

		err = json.Unmarshal(rawBody, &body)
		if err != nil {
			return nil, providererrors.NewErrRequestWithResponse(req, res)
		}

		return body.Data, nil

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
