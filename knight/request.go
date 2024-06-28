package knight

import (
	"errors"
)

var errMissingValue = errors.New("request: key missing in params")

type Request struct {
	From   string
	Hit    string // like url
	Params map[string]string
}

func (r *Request) SetHit(h string) {
	r.Hit = h
}

func (r *Request) SetParam(key, val string) {
	if r.Params == nil {
		r.Params = make(map[string]string)
	}
	r.Params[key] = val
}

func (r *Request) GetParam(key string) (string, error) {
	if val, ok := r.Params[key]; ok {
		return val, nil
	}
	return "", errMissingValue
}
