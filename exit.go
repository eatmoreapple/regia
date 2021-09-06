// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

type Exit interface{ Exit(context *Context) }

type exit struct{}

// Exit Do nothing
func (e exit) Exit(*Context) {}
