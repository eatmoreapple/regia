// Copyright 2022 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import "net/url"

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

func (ps Params) Get(key string) string {
	return ps.byName(key)
}

// ToURLValues converts a Params to an url.Values
// This is useful for building a URL query string
func (ps Params) ToURLValues() url.Values {
	values := url.Values{}
	for _, p := range ps {
		values.Add(p.Key, p.Value)
	}
	return values
}
