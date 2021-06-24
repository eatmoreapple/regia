package regia

import (
	"strconv"
	"sync"
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

func (v Value) Int8(def ...int8) (int8, error) {
	i, err := strconv.ParseInt(string(v), 10, 8)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return int8(i), err
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

func (v Value) Int32(def ...int32) (int32, error) {
	i, err := strconv.ParseInt(string(v), 10, 32)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return int32(i), err
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

func (v Value) Uint(def ...uint) (uint, error) {
	i, err := strconv.ParseUint(string(v), 10, 0)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return uint(i), err
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

func (v Value) Uint16(def ...uint16) (uint16, error) {
	i, err := strconv.ParseUint(string(v), 10, 16)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return uint16(i), err
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

func (v Value) Uint64(def ...uint64) (uint64, error) {
	i, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return i, err
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

func (v Value) Float64(def ...float64) (float64, error) {
	i, err := strconv.ParseFloat(string(v), 64)
	if err != nil || v.IsEmpty() {
		if len(def) > 0 {
			return def[0], err
		}
	}
	return i, err
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

func (v Value) Duration(def ...int64) (time.Duration, error) {
	i, err := v.Int64(def...)
	return time.Duration(i), err
}

func (v Value) Unmarshal(f Unmarshaler, dst interface{}) error {
	return f.Unmarshal([]byte(v), dst)
}

type Values []Value

type Warehouse interface {
	Set(key string, value interface{})
	Get(key string) (value interface{}, exist bool)
}

var _ Warehouse = new(SyncMap)

type SyncMap struct {
	item map[string]interface{}
	mu   sync.RWMutex
}

func (d *SyncMap) Set(key string, value interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.item == nil {
		d.item = make(map[string]interface{})
	}
	d.item[key] = value
}

func (d *SyncMap) Get(key string) (value interface{}, exist bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	value, exist = d.item[key]
	return
}

func (d *SyncMap) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.item = nil
}

func (d *SyncMap) Reset() {
	d.Clear()
}

func (d *SyncMap) GetString(key string, def ...string) string {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(string); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (d *SyncMap) GetInt(key string, def ...int) int {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(int); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetInt8(key string, def ...int8) int8 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(int8); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetInt16(key string, def ...int16) int16 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(int16); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetInt32(key string, def ...int32) int32 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(int32); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetInt64(key string, def ...int64) int64 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(int64); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetUint(key string, def ...uint) uint {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(uint); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetUint8(key string, def ...uint8) uint8 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(uint8); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetUint16(key string, def ...uint16) uint16 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(uint16); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetUint32(key string, def ...uint32) uint32 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(uint32); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetUint64(key string, def ...uint64) uint64 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(uint64); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetFloat32(key string, def ...float32) float32 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(float32); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetFloat64(key string, def ...float64) float64 {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(float64); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return 0
}

func (d *SyncMap) GetBool(key string, def ...bool) bool {
	value, exist := d.Get(key)
	if !exist {
		if len(def) > 0 {
			return def[0]
		}
	}
	if i, ok := value.(bool); ok {
		return i
	} else if len(def) > 0 {
		return def[0]
	}
	return false
}

const (
	author = "多吃点苹果"
	wechat = "EatMoreApple"
)
