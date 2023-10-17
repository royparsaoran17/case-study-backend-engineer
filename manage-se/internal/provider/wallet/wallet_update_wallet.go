package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"manage-se/internal/presentations"
	"manage-se/internal/provider/providererrors"
	"manage-se/pkg/httpx"
	"net/http"
)

func (c *client) UpdateWallet(ctx context.Context, walletID string, input presentations.WalletUpdate) error {
	urlEndpoint := c.endpoint(fmt.Sprintf("/internal/v1/wallets/%s", walletID))

	var request bytes.Buffer
	err := json.NewEncoder(&request).Encode(input)
	if err != nil {
		return errors.Wrap(err, "new encoder encode")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, urlEndpoint, &request)
	if err != nil {
		return errors.Wrap(err, "new request with context")
	}

	req.Header.Set(httpx.ContentType, httpx.MediaTypeJSON)

	res, err := c.dep.HttpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
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
			return providererrors.NewErrRequestWithResponse(req, res)
		}

		return nil

	default:
		bodyErr := providererrors.Error{}
		err := json.Unmarshal(rawBody, &bodyErr)
		if err != nil {
			return providererrors.NewErrRequestWithResponse(req, res)
		}

		bodyErr.Code = res.StatusCode
		return bodyErr

	}
}
