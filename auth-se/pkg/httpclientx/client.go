package httpclientx

import (
	"auth-se/internal/consts"
	"auth-se/internal/presentations"
	"auth-se/pkg/tracer"
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
)

type HttpClientx struct {
	httpClient *http.Client
}

func NewHttpClientx() *HttpClientx {
	return &HttpClientx{httpClient: http.DefaultClient}
}

func (h *HttpClientx) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	tracer.AddSpanTag(ctx,
		tracer.NewSpanTag("http.request.headers.*", req.Header),
		tracer.NewSpanTag("http.url", req.URL.Path),
		tracer.NewSpanTag("http.method", req.Method),
	)

	if req.Body != nil {
		// Re-usable request body for logging
		requestBody, _ := io.ReadAll(req.Body)
		req.Body.Close() // must close
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		tracer.AddSpanTag(ctx,
			tracer.NewSpanTag("http.request.body", string(requestBody)),
		)
	}

	//Add X-Request-ID to the request Header
	state, valid := ctx.Value(consts.CtxRequestState).(presentations.RequestState)
	if valid {
		req.Header.Set(consts.HeaderXRequestID, state.ID)
	}

	// Attempt to join a trace by getting trace context from the headers.
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	span := opentracing.SpanFromContext(ctx)

	// Transmit the span's TraceContext as HTTP headers on our
	// outbound request.
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		carrier)

	res, err := h.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Re-usable response body for logging
	responseBody, _ := io.ReadAll(res.Body)
	res.Body.Close() // must close
	res.Body = io.NopCloser(bytes.NewBuffer(responseBody))

	tracer.AddSpanTag(ctx,
		tracer.NewSpanTag("http.status_code", strconv.Itoa(res.StatusCode)),
		tracer.NewSpanTag("http.response.headers.*", res.Header),
		tracer.NewSpanTag("http.response.body", string(responseBody)),
	)

	return res, nil
}
