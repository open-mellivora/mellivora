package core

import (
	"encoding/json"
	"errors"
	"reflect"
)

type setter struct {
	write map[string]interface{}
	read  map[string]json.RawMessage
}

func newSetter() setter {
	return setter{write: make(map[string]interface{})}
}

func (c *setter) MarshalText() ([]byte, error) {
	return json.Marshal(c.write)
}

func (c *setter) UnmarshalText(text []byte) error {
	err := json.Unmarshal(text, &c.read)
	return err
}

func (c *setter) Set(k string, v interface{}) {
	c.write[k] = v
}

func (c *setter) Value(k string, v interface{}) error {
	if c.write != nil {
		wv, ok := c.write[k]
		if !ok {
			return nil
		}

		valueOf := reflect.ValueOf(v)
		if valueOf.Kind() != reflect.Ptr || valueOf.IsNil() {
			return errors.New("value not ptr")
		}

		if !valueOf.Elem().CanSet() {
			return errors.New("value cant set")
		}

		valueOf.Elem().Set(reflect.ValueOf(wv))
		return nil
	}
	var bs []byte
	var has bool
	if bs, has = c.read[k]; !has {
		return nil
	}

	return json.Unmarshal(bs, v)
}

func (c *setter) MustValue(k string, v interface{}) {
	if err := c.Value(k, v); err != nil {
		panic(err)
	}
}

// SetDontFilter sets `depth`.
func (c *setter) SetDontFilter(dontFilter bool) {
	c.Set(dontFilterKey, dontFilter)
}

// GetDontFilter returns `depth`.
func (c *setter) GetDontFilter() bool {
	var value bool
	c.MustValue(dontFilterKey, &value)
	return value
}

// SetDepth sets `depth`.
func (c *setter) SetDepth(depth int64) {
	c.Set(depthKey, depth)
}

// GetDepth returns `depth`.
func (c *setter) GetDepth() int64 {
	var value int64
	c.MustValue(depthKey, &value)
	return value
}
