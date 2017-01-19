package jx

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/mitchellh/mapstructure"
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

func (j *JX) Value() interface{} {
	return j.x
}

func (j *JX) Get(key string) *JX {
	if x, ok := j.x.(map[string]interface{}); ok {
		return &JX{x[key]}
	}
	return &JX{new(interface{})}
}

func (j *JX) Bool(key ...string) (r bool) {
	if len(key) > 0 {
		return j.Get(key[0]).Bool()
	}
	if v, ok := j.x.(bool); ok {
		return v
	}
	return r
}

func (j *JX) String(key ...string) (r string) {
	if len(key) > 0 {
		return j.Get(key[0]).String()
	}
	if v, ok := j.x.(string); ok {
		return v
	}
	return r
}

func (j *JX) Int(key ...string) (r int) {
	if len(key) > 0 {
		return j.Get(key[0]).Int()
	}
	if v, ok := j.x.(int); ok {
		return v
	}
	return r
}

func (j *JX) Float(key ...string) (r float64) {
	if len(key) > 0 {
		return j.Get(key[0]).Float()
	}
	if v, ok := j.x.(float64); ok {
		return v
	}
	return r
}

func (j *JX) Map(key ...string) (r map[string]interface{}) {
	if len(key) > 0 {
		return j.Get(key[0]).Map()
	}
	if v, ok := j.x.(map[string]interface{}); ok {
		return v
	}
	return r
}

func (j *JX) Slice(key ...string) (r []interface{}) {
	if len(key) > 0 {
		return j.Get(key[0]).Slice()
	}
	if v, ok := j.x.([]interface{}); ok {
		return v
	}
	return r
}

type Item struct {
	Key string
	*JX
}

func (j *JX) MapIter(key ...string) chan Item {
	if len(key) > 0 {
		return j.Get(key[0]).MapIter()
	}
	ch := make(chan Item)
	go func() {
		for k, v := range j.Map() {
			ch <- Item{k, &JX{v}}
		}
		close(ch)
	}()
	return ch
}

func (j *JX) SliceIter(key ...string) chan *JX {
	if len(key) > 0 {
		return j.Get(key[0]).SliceIter()
	}
	ch := make(chan *JX)
	go func() {
		for _, v := range j.Slice() {
			ch <- &JX{v}
		}
		close(ch)
	}()
	return ch
}

func (j *JX) Decode(into interface{}, key ...string) error {
	if len(key) > 0 {
		return j.Get(key[0]).Decode(into)
	}
	return mapstructure.Decode(j.Map(), into)
}
