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

func (j *JX) Bool() (r bool) {
	if v, ok := j.x.(bool); ok {
		return v
	}
	return r
}

func (j *JX) String() (r string) {
	if v, ok := j.x.(string); ok {
		return v
	}
	return r
}

func (j *JX) Int() (r int) {
	if v, ok := j.x.(int); ok {
		return v
	}
	return r
}

func (j *JX) Float() (r float64) {
	if v, ok := j.x.(float64); ok {
		return v
	}
	return r
}

func (j *JX) Map() (r map[string]interface{}) {
	if v, ok := j.x.(map[string]interface{}); ok {
		return v
	}
	return r
}

func (j *JX) Slice() (r []interface{}) {
	if v, ok := j.x.([]interface{}); ok {
		return v
	}
	return r
}

type Item struct {
	Key string
	*JX
}

func (j *JX) MapIter() chan Item {
	ch := make(chan Item)
	go func() {
		for k, v := range j.Map() {
			ch <- Item{k, &JX{v}}
		}
		close(ch)
	}()
	return ch
}

func (j *JX) SliceIter() chan *JX {
	ch := make(chan *JX)
	go func() {
		for _, v := range j.Slice() {
			ch <- &JX{v}
		}
		close(ch)
	}()
	return ch
}

func (j *JX) Decode(into interface{}) error {
	return mapstructure.Decode(j.Map(), into)
}
