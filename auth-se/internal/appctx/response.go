// Package appctx
package appctx

import (
	"auth-se/internal/consts"
	"auth-se/pkg/msgx"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"sync"
)

var (
	rsp    *Response
	oneRsp sync.Once
)

// StatusAPIXXX are statuses for API Response.
const (
	StatusAPISuccess = "SUCCESS"
	StatusAPIError   = "ERROR"
	StatusAPIFailure = "FAILURE"
)

// Response presentation contract object
type Response struct {
	Status  string      `json:"status,omitempty"`
	Entity  string      `json:"entity,omitempty"`
	State   string      `json:"state,omitempty"`
	Code    int         `json:"code"`
	Message interface{} `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	lang    string      `json:"-"`
	Meta    interface{} `json:"meta,omitempty"`
	msgKey  string
	Custom  *ResponseCustom `json:"-"`
	ctx     context.Context
}

type ResponseCustom struct {
	ContentType string
	Content     []byte
}

// MetaData represent meta data response for multi data
type MetaData struct {
	Page       uint64 `json:"page"`
	Limit      uint64 `json:"limit"`
	TotalPage  uint64 `json:"total_page"`
	TotalCount uint64 `json:"total_count"`
}

// GetMessage method to transform response name var to message detail
func (r *Response) GetMessage() string {
	return msgx.Get(r.msgKey, r.lang).Text()
}

// Generate setter message
func (r *Response) Generate() *Response {
	if r.lang == "" {
		r.lang = consts.LangDefault
	}
	msg := msgx.Get(r.msgKey, r.lang)
	if r.Message == nil {
		r.Message = msg.Text()
	}

	if r.Code == 0 {
		r.Code = msg.Status()
	}

	return r
}

// WithCode setter response var name
func (r *Response) WithCode(c int) *Response {
	r.Code = c

	if r.Code/100 == 5 {
		r.Status = StatusAPIFailure
	} else if r.Code/100 == 4 {
		r.Status = StatusAPIError
	} else {
		r.Status = StatusAPISuccess
	}

	caser := cases.Title(language.English)
	state := caser.String(r.Status)
	r.State = fmt.Sprintf("%s%s", r.Entity, state)

	return r
}

// WithData setter data response
func (r *Response) WithData(v interface{}) *Response {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *Response) WithError(v interface{}) *Response {
	r.Errors = v
	return r
}

func (r *Response) WithMsgKey(v string) *Response {
	r.msgKey = v
	return r
}

// WithMeta setter meta data response
func (r *Response) WithMeta(v interface{}) *Response {
	r.Meta = v
	return r
}

// WithLang setter language response
func (r *Response) WithLang(v string) *Response {
	r.lang = v
	return r
}

// WithMessage setter custom message response
func (r *Response) WithMessage(v interface{}) *Response {
	if v != nil {
		r.Message = v
	}

	return r
}

// WithState setter response entity name and state
func (r *Response) WithState(entity string) *Response {
	r.Entity = entity

	return r
}

func (r *Response) Byte() []byte {
	if r.Code == 0 || r.Message == nil {
		r.Generate()
	}

	b, _ := json.Marshal(r)
	return b
}

func (r *Response) WithCustomResponse(custom ResponseCustom) *Response {
	r.Custom = &custom
	return r
}

func (r *Response) WithContext(ctx context.Context) *Response {
	r.ctx = ctx
	return r
}

func (r *Response) Context() context.Context {
	return r.ctx
}

// NewResponse initialize response
func NewResponse() *Response {
	oneRsp.Do(func() {
		rsp = &Response{}
	})

	// clone response
	x := *rsp

	return &x
}
