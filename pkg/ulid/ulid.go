package ulid

/*
 * Copyright 2021 OpsMx, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	cryptorand "crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// Context holds the state needed to generate a ULID using random values.  This
// is a locked structure, so should not be used by a lot of threads if IDs are
// generated at a high rate.
type Context struct {
	entropy *ulid.MonotonicEntropy
}

// NewContext returns the context needed for subsequent calls.
func NewContext() *Context {
	entropy := ulid.Monotonic(cryptorand.Reader, 0)
	return &Context{entropy: entropy}
}

// Ulid - return a new ULID as a string.
func (ctx *Context) Ulid() string {
	t := time.Now().Unix()
	return ulid.MustNew(uint64(t), ctx.entropy).String()
}
