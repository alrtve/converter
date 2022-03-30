package converters

import (
	"encoding/json"
	"io"
	"io/ioutil"
)
import "gopkg.in/yaml.v2"

type YmlToJsonConverter struct {
}

func NewYmlToJsonConverter() *YmlToJsonConverter {
	return &YmlToJsonConverter{}
}

func (c *YmlToJsonConverter) Convert(w io.Writer, r io.Reader) error {
	var tmp map[interface{}]interface{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	tmpN := c.convert(tmp)
	if data, err = json.Marshal(tmpN); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (c *YmlToJsonConverter) convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = c.convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = c.convert(v)
		}
	}
	return i
}
