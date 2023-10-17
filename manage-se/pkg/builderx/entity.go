package builderx

type (
	KeyValue struct {
		Key       string
		Value     any
		IsPrimary bool
	}

	QueryWhere struct {
		Columns []string
		Values  []any
		Limit   int64
		Page    int64
		Query   string
	}
)
