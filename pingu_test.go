package pingu_test

import (
	"testing"

	"github.com/dbadoy/pingu"
)

func TestNewPinguAddr(t *testing.T) {
	type td struct {
		got    string
		expect pingu.Config
		err    string
	}
	tdl := []td{
		{got: "127.0.0.1:1723", expect: pingu.Config{}, err: ""},
		{got: "127.0.0.1:", expect: pingu.Config{}, err: "no port"},
		{got: ":3821", expect: pingu.Config{}, err: "no IP"},
		{got: "", expect: pingu.Config{}, err: "not an ip:port"},
		{got: "hello", expect: pingu.Config{}, err: "not an ip:port"},
	}

	for _, td := range tdl {
		_, err := pingu.NewPingu(td.got, nil)
		if err != nil {
			if err.Error() != td.err {
				t.Fatalf("NewPingu failure got: %v, want: %v", err.Error(), td.err)
			}
			continue
		}
	}
}
