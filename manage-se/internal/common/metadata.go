package common

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ErrInvalidMetadata is an error when metadata is invalid.
// This error usually returned by the implementation of Filter interface.
var ErrInvalidMetadata = errors.New("invalid metadata")

var DefaultMetaData = MetadataFromURL(url.Values{})

// Metadata represents a metadata for HTTP API.
type Metadata struct {
	Pagination
	Filtering
	*DateRange `json:"date_range,omitempty"`
}

// MetadataFromURL gets metadata from the given request url.
func MetadataFromURL(u url.Values) Metadata {
	return Metadata{
		Pagination: PaginationFromURL(u),
		Filtering:  FilterFromURL(u),
	}
}

type MetaDataOpt struct {
	dateRange *DateRange
}

type MetaDataOpts func(opt *MetaDataOpt)

func MetaDataWithDateRange(u url.Values, field string, startQuery, endQuery string) (MetaDataOpts, error) {
	dr, err := DateRangeFromURL(u, field, startQuery, endQuery)
	if err != nil {
		return nil, err
	}

	return func(opt *MetaDataOpt) {
		opt.dateRange = dr
	}, nil
}

// DefaultPerPage is a default value for per_page query params.
const DefaultPerPage = 10

// Pagination is a meta data for pagination.
type Pagination struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
	Total   int `json:"total"`
}

// PaginationFromURL gets pagination meta from request URL.
func PaginationFromURL(u url.Values) Pagination {
	p := Pagination{
		PerPage: DefaultPerPage,
		Page:    1,
	}

	pps := u.Get("per_page")
	if v, err := strconv.Atoi(pps); err == nil {
		if v <= 0 {
			v = DefaultPerPage
		}

		p.PerPage = v
	}

	ps := u.Get("page")
	if v, err := strconv.Atoi(ps); err == nil {
		if v < 1 {
			v = 1
		}

		p.Page = v
	}

	return p
}

// SortXXX are default values for order_type query params.
const (
	SortAscending  = "asc"
	SortDescending = "desc"
)

// Filtering represents a filterable fields.
type Filtering struct {
	OrderBy   string `json:"order_by"`
	OrderType string `json:"order_type"`
	Search    string `json:"search,omitempty"`
	SearchBy  string `json:"search_by,omitempty"`
}

// FilterFromURL gets filter values from query params.
func FilterFromURL(u url.Values) Filtering {
	f := Filtering{
		OrderBy:   "created_at",
		OrderType: SortAscending,
	}

	ob := u.Get("order_by")
	if len(ob) != 0 {
		f.OrderBy = strings.ToLower(strings.ToLower(ob))
	}

	ot := u.Get("order_type")
	if len(ot) != 0 {
		ot = strings.TrimSpace(strings.ToLower(ot))
		if ot == SortDescending {
			f.OrderType = SortDescending
		}
	}

	search := strings.TrimSpace(u.Get("search"))
	if len(search) == 0 {
		search = strings.TrimSpace(u.Get("keyword"))
	}

	if len(search) != 0 {
		f.Search = search
	}

	searchBy := strings.TrimSpace(u.Get("search_by"))
	if len(searchBy) != 0 {
		f.SearchBy = searchBy
	}

	return f
}

type DateRange struct {
	Field string    `json:"field"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func DateRangeFromURL(u url.Values, field string, startQuery, endQuery string) (*DateRange, error) {
	ts := u.Get(startQuery)
	te := u.Get(endQuery)
	if len(ts) == 0 || len(te) == 0 {
		return nil, nil
	}

	dr := DateRange{
		Field: "created_at",
		Start: time.Time{},
		End:   time.Time{},
	}

	if v := u.Get(field); len(v) != 0 {
		dr.Field = strings.TrimSpace(strings.ToLower(v))
	}

	t, err := time.Parse("2006-01-02", ts)
	if err != nil {
		return nil, ErrInvalidMetadata
	}

	dr.Start = t

	t, err = time.Parse("2006-01-02", te)
	if err != nil {
		return nil, ErrInvalidMetadata
	}

	dr.End = t

	return &dr, nil
}

// Filter knows how to validate filterable fields.
// This Filter usually implemented by Repository.
type Filter interface {
	// Sortable returns true if a given field is allowed for sorting.
	Sortable(field string) bool
}

type Query struct {
	OrderBy        string
	OrderDirection string
	Search         string
	SearchBy       string
	Limit          int
	Offset         int
	DateFrom       sql.NullTime
	DateEnd        sql.NullTime
}

func ParamFromMetadata(metadata *Metadata, filter Filter) (*Query, error) {
	if !filter.Sortable(metadata.OrderBy) {
		return nil, ErrInvalidMetadata
	}

	var form, end sql.NullTime
	if metadata.DateRange != nil {
		if !filter.Sortable(metadata.DateRange.Field) {
			return nil, ErrInvalidMetadata
		}

		form = sql.NullTime{
			Time:  BeginOfDay(metadata.DateRange.Start),
			Valid: !metadata.DateRange.Start.IsZero(),
		}

		end = sql.NullTime{
			Time:  BeginOfNextDay(metadata.DateRange.End),
			Valid: !metadata.DateRange.End.IsZero(),
		}
	}

	limit := metadata.PerPage
	offset := (metadata.Page - 1) * limit
	search := "%" + strings.ToLower(metadata.Search) + "%"

	q := Query{
		OrderBy:        metadata.OrderBy,
		OrderDirection: metadata.OrderType,
		Search:         search,
		Limit:          limit,
		Offset:         offset,
		DateFrom:       form,
		DateEnd:        end,
		SearchBy:       strings.ToLower(metadata.SearchBy),
	}

	return &q, nil
}

func (m Metadata) ToURLValues() url.Values {
	v := url.Values{}

	v.Set("page", fmt.Sprintf("%d", m.Page))
	v.Set("per_page", fmt.Sprintf("%d", m.PerPage))
	v.Set("order_by", m.OrderBy)
	v.Set("order_type", m.OrderType)
	v.Set("search_by", m.SearchBy)
	v.Set("search", m.Search)
	if m.DateRange != nil {
		v.Set("created_from", m.DateRange.Start.Format("2006-01-02"))
		v.Set("created_until", m.DateRange.End.Format("2006-01-02"))
	}

	return v
}
