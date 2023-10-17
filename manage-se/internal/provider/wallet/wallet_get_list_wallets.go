package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"manage-se/internal/common"
	"manage-se/internal/provider/providererrors"
	"net/http"
)

func (c *client) GetListWallets(ctx context.Context, meta *common.Metadata) ([]Wallet, error) {
	urlEndpoint := c.endpoint("/internal/v1/wallets")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlEndpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request failed")
	}

	q := req.URL.Query()

	q.Set("page", fmt.Sprintf("%d", meta.Page))
	q.Set("per_page", fmt.Sprintf("%d", meta.PerPage))
	q.Set("search", fmt.Sprintf("%s", meta.Search))
	q.Set("search_by", fmt.Sprintf("%s", meta.SearchBy))
	q.Set("order_type", fmt.Sprintf("%s", meta.OrderType))
	q.Set("order_by", fmt.Sprintf("%s", meta.OrderBy))
	if meta.DateRange != nil {
		q.Set("created_from", meta.DateRange.Start.Format("2006-01-02"))
		q.Set("created_until", meta.DateRange.End.Format("2006-01-02"))
	}

	req.URL.RawQuery = q.Encode()

	res, err := c.dep.HttpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("doing http request to %s", req.URL))
	}

	// Re-usable response body for logging
	rawBody, _ := io.ReadAll(res.Body)
	res.Body.Close() // must close
	res.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	switch res.StatusCode {
	case http.StatusOK:
		body := struct {
			Data []Wallet `json:"data"`
			Meta Metadata `json:"meta"`
		}{}

		err = json.Unmarshal(rawBody, &body)
		if err != nil {
			return nil, providererrors.NewErrRequestWithResponse(req, res)
		}

		meta.Total = body.Meta.Total
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
