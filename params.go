package dsnparser

import (
	"net/url"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// DecodeParams decode params url into a struct
func (d *DSN) DecodeParams(value interface{}) error {
	var err error
	urlval := url.Values{}
	for k, v := range d.params {
		urlval.Add(k, v)
	}
	err = decoder.Decode(value, urlval)
	if err != nil {
		return err
	}

	return nil
}
