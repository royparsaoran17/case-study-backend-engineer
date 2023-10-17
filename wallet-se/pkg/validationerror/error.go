package validationerror

import (
	"encoding/json"
)

type Error map[string][]string

func (e Error) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

func NewError() Error {
	return make(Error)
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (e Error) Get(key string) string {
	if e == nil {
		return ""
	}
	vs := e[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (e Error) Set(key, value string) {
	e[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (e Error) Add(key, value string) {
	e[key] = append(e[key], value)
}

// Del deletes the values associated with key.
func (e Error) Del(key string) {
	delete(e, key)
}

// Has checks whether a given key is set.
func (e Error) Has(key string) bool {
	_, ok := e[key]
	return ok
}
