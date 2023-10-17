package common

import (
	"math"

	"auth-se/internal/consts"
)

// LimitDefaultValue set default value limit
func LimitDefaultValue(origin int64) int64 {
	if origin < 1 {
		return consts.PagingLimitDefaultValue
	}

	if origin > consts.PagingMaxLimit {
		return consts.PagingMaxLimit
	}

	return origin
}

// PageDefaultValue set default value page
func PageDefaultValue(origin int64) int64 {
	if origin < 1 {
		return consts.PagingPageDefaultValue
	}

	return origin
}

// PageCalculate calculate total page from count
func PageCalculate(count int64, limit int64) int64 {
	if count <= limit {
		return 1
	}

	return int64(math.Ceil(float64(count) / float64((limit))))
}

// OffsetDefaultValue set default offset
func OffsetDefaultValue(page int64, limit int64) int64 {
	if page < 1 {
		return consts.PagingOffsetDefaultValue
	}

	return (page - 1) * limit
}

// PageToOffset calculate
func PageToOffset(limit, page int64) int64 {
	if page <= 0 {
		return 0
	}

	return (page - 1) * limit
}
