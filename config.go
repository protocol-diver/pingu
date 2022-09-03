// Copyright (c) 2022, Seungbae Yu <dbadoy4874@gmail.com>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pingu

type Config struct {
	RecvBufferSize int
	Verbose        bool
}

func (c *Config) Default() {
	c.RecvBufferSize = 256
	c.Verbose = false
}
