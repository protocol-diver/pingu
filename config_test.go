package pingu_test

import (
	"testing"

	"github.com/protocol-diver/pingu"
)

func TestConfigDeafult(t *testing.T) {
	type td struct {
		got    pingu.Config
		expect pingu.Config
	}

	defaultBufferSize := pingu.DefultRecvBufferSize
	tdl := []td{
		{got: pingu.Config{RecvBufferSize: 5, Verbose: true}, expect: pingu.Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: pingu.Config{RecvBufferSize: 3, Verbose: false}, expect: pingu.Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: pingu.Config{RecvBufferSize: 80}, expect: pingu.Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: pingu.Config{Verbose: true}, expect: pingu.Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: pingu.Config{Verbose: false}, expect: pingu.Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: pingu.Config{}, expect: pingu.Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
	}

	for _, td := range tdl {
		td.got.Default()
		if !compare(td.got, td.expect) {
			t.Fatalf("Config.Deafult failure got: %v, want: %v", td.got, td.expect)
		}
	}
}

func compare(a pingu.Config, b pingu.Config) bool {
	if a.RecvBufferSize != b.RecvBufferSize {
		return false
	}
	if a.Verbose != b.Verbose {
		return false
	}
	return true
}
