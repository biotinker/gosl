// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnd

import (
	"github.com/biotinker/gosl/chk"
	"github.com/biotinker/gosl/io"
)

func init() {
	io.Verbose = false
}

func verbose() {
	io.Verbose = true
	chk.Verbose = true
}
