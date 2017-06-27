package jx

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type JX struct {
	x interface{}
}

func FromReader(reader io.Reader) (*JX, error) {
	j := new(JX)
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	err := decoder.Decode(&j.x)
	return j, err
}

func FromBytes(data []byte) (*JX, error) {
	return FromReader(bytes.NewBuffer(data))
}

func FromFile(path string) (*JX, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return FromReader(file)
}

func empty() *JX {
	return &JX{new(interface{})}
}

func (j *JX) Index(i int) *JX {
	if v, ok := j.x.([]interface{}); ok && i < len(v) {
		return &JX{v[i]}
	}
	return empty()
}

func (j *JX) Key(key string) *JX {
	if v, ok := j.x.(map[string]interface{}); ok {
		return &JX{v[key]}
	}
	return empty()
}

func (j *JX) Get(path ...interface{}) *JX {
	var current = j
	for _, key := range path {
		switch actual := key.(type) {
		case int:
			current = current.Index(actual)
		case string:
			current = current.Key(actual)
		default:
			return empty()
		}
	}
	return current
}

func (j *JX) Value() interface{} {
	return j.x
}

func (j *JX) Int() (int64, bool) {
	if number, ok := j.x.(json.Number); ok {
		value, err := number.Int64()
		return value, err == nil
	}
	return 0, false
}

func (j *JX) Float() (float64, bool) {
	if number, ok := j.x.(json.Number); ok {
		value, err := number.Float64()
		return value, err == nil
	}
	return 0, false
}

func (j *JX) Bool() (bool, bool) {
	value, ok := j.x.(bool)
	return value, ok
}

func (j *JX) String() (string, bool) {
	value, ok := j.x.(string)
	return value, ok
}
