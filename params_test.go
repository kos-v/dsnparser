package dsnparser

import (
	"fmt"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

type ParamOption struct {
	Tblsprefix string
	Charset    string
	ParseTime  bool
	Loc        string
	Topic      []string
	Version    string
}

func TestDecodeMapToStruct(t *testing.T) {

	var err error
	var testcases = []struct {
		dsn       string
		param     ParamOption
		wantError bool
	}{
		{
			dsn: "mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_",
			param: ParamOption{
				Tblsprefix: "fs_",
			},
			wantError: false,
		},
		{
			dsn: "mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_&charset=utf8mb4&parseTime=True&loc=Local",
			param: ParamOption{
				Tblsprefix: "fs_",
				Charset:    "utf8mb4",
				ParseTime:  true,
				Loc:        "Local",
			},
			wantError: false,
		},
		{
			dsn: "kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog&version=1.1.1",
			param: ParamOption{
				Topic:   []string{"vsulblog"},
				Version: "1.1.1",
			},
			wantError: false,
		},
	}

	for _, tt := range testcases {
		var res = ParamOption{}
		dsn := Parse(tt.dsn)
		err = dsn.DecodeParams(&res)
		fmt.Println(err)
		assert.Equal(t, err != nil, tt.wantError)
		if err != nil {
			continue
		}
		assert.Equal(t, res, tt.param)
	}

}
