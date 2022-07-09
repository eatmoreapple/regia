// Copyright 2022 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) byName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

func (ps Params) Get(key string) Value {
	v := ps.byName(key)
	return Value(v)
}
