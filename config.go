package pingu

type Config struct {
	RecvBufferSize int
	Verbose        bool
}

func (c *Config) Default() {
	c.RecvBufferSize = 512
	c.Verbose = false
}
