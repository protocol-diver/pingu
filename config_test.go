package pingu

import (
	"testing"
)

func TestConfigDeafult(t *testing.T) {
	type td struct {
		got    Config
		expect Config
	}

	defaultBufferSize := DefultRecvBufferSize
	tdl := []td{
		{got: Config{RecvBufferSize: 5, Verbose: true}, expect: Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: Config{RecvBufferSize: 3, Verbose: false}, expect: Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: Config{RecvBufferSize: 80}, expect: Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: Config{Verbose: true}, expect: Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: Config{Verbose: false}, expect: Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
		{got: Config{}, expect: Config{RecvBufferSize: defaultBufferSize, Verbose: false}},
	}

	for _, td := range tdl {
		td.got.Default()
		if !compare(td.got, td.expect) {
			t.Fatalf("Config.Deafult failure got: %v, want: %v", td.got, td.expect)
		}
	}
}

func compare(a Config, b Config) bool {
	if a.RecvBufferSize != b.RecvBufferSize {
		return false
	}
	if a.Verbose != b.Verbose {
		return false
	}
	return true
}
