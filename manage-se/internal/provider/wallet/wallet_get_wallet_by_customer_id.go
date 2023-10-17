package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"manage-se/internal/provider/providererrors"
	"manage-se/pkg/httpx"
	"net/http"
)

func (c *client) GetWalletByCustomerID(ctx context.Context, customerID string) (*WalletDetail, error) {
	urlEndpoint := c.endpoint(fmt.Sprintf("/internal/v1/wallets/owned/%s", customerID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEndpoint, nil)
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
			Data WalletDetail `json:"data"`
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
