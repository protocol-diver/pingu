package pingu

type Config struct {
	RecvBufferSize int
	Verbose        bool
}

func (c *Config) Default() {
	c.RecvBufferSize = 256
	c.Verbose = false
}
