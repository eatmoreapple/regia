// Copyright 2021 eatmoreapple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"strconv"
	"time"
)

type Unmarshaler interface {
	Unmarshal(data []byte, v interface{}) error
}

type Value string

func (v Value) IsEmpty() bool {
	return v == ""
}

func (v Value) Text(def ...string) string {
	if v.IsEmpty() && len(def) > 0 {
		return def[0]
	}
	return string(v)
}

func (v Value) Int(def ...int) (int, error) {
	i, err := strconv.ParseInt(string(v), 10, 0)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return int(i), err
}

func (v Value) MustInt(def ...int) int {
	i, _ := v.Int(def...)
	return i
}

func (v Value) Int8(def ...int8) (int8, error) {
	i, err := strconv.ParseInt(string(v), 10, 8)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return int8(i), err
}

func (v Value) MustInt8(def ...int8) int8 {
	i, _ := v.Int8(def...)
	return i
}

func (v Value) Int16(def ...int16) (int16, error) {
	i, err := strconv.ParseInt(string(v), 10, 16)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return int16(i), err
}

func (v Value) MustInt16(def ...int16) int16 {
	i, _ := v.Int16(def...)
	return i
}

func (v Value) Int32(def ...int32) (int32, error) {
	i, err := strconv.ParseInt(string(v), 10, 32)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return int32(i), err
}

func (v Value) MustInt32(def ...int32) int32 {
	i, _ := v.Int32(def...)
	return i
}

func (v Value) Int64(def ...int64) (int64, error) {
	i, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return i, err
}

func (v Value) MustInt64(def ...int64) int64 {
	i, _ := v.Int64(def...)
	return i
}

func (v Value) Uint(def ...uint) (uint, error) {
	i, err := strconv.ParseUint(string(v), 10, 0)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return uint(i), err
}

func (v Value) MustUint(def ...uint) uint {
	i, _ := v.Uint(def...)
	return i
}

func (v Value) Uint8(def ...uint8) (uint8, error) {
	i, err := strconv.ParseUint(string(v), 10, 8)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return uint8(i), err
}

func (v Value) MustUint8(def ...uint8) uint8 {
	i, _ := v.Uint8(def...)
	return i
}

func (v Value) Uint16(def ...uint16) (uint16, error) {
	i, err := strconv.ParseUint(string(v), 10, 16)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return uint16(i), err
}

func (v Value) MustUint16(def ...uint16) uint16 {
	i, _ := v.Uint16(def...)
	return i
}

func (v Value) Uint32(def ...uint32) (uint32, error) {
	i, err := strconv.ParseUint(string(v), 10, 32)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return uint32(i), err
}

func (v Value) MustUint32(def ...uint32) uint32 {
	i, _ := v.Uint32(def...)
	return i
}

func (v Value) Uint64(def ...uint64) (uint64, error) {
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return i, err
}

func (v Value) MustUint64(def ...uint64) uint64 {
	i, _ := v.Uint64(def...)
	return i
}

func (v Value) Float32(def ...float32) (float32, error) {
	i, err := strconv.ParseFloat(string(v), 32)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return float32(i), err
}

func (v Value) MustFloat32(def ...float32) float32 {
	i, _ := v.Float32(def...)
	return i
}

func (v Value) Float64(def ...float64) (float64, error) {
	i, err := strconv.ParseFloat(string(v), 64)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return i, err
}

func (v Value) MustFloat64(def ...float64) float64 {
	i, _ := v.Float64(def...)
	return i
}

func (v Value) Bool(def ...bool) (bool, error) {
	if v.IsEmpty() {
		return false, nil
	}
	i, err := strconv.ParseBool(string(v))
	if err != nil {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return i, err
}

func (v Value) MustBool(def ...bool) bool {
	i, _ := v.Bool(def...)
	return i
}

func (v Value) Duration(def ...int64) (time.Duration, error) {
	i, err := v.Int64(def...)
	return time.Duration(i), err
}

func (v Value) Unmarshal(f Unmarshaler, dst interface{}) error {
	return f.Unmarshal([]byte(v), dst)
}

type Values []Value

func NewValues(values []string) Values {
	var v = make(Values, len(values))
	for _, item := range values {
		v = append(v, Value(item))
	}
	return v
}

type Warehouse interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, exist bool)
}

var _ Warehouse = warehouse(nil)

type warehouse map[interface{}]interface{}

func (w warehouse) Set(key, value interface{}) {
	w[key] = value
}

func (w warehouse) Get(key interface{}) (interface{}, bool) {
	value, exist := w[key]
	return value, exist
}
