package converters

import (
	"encoding/json"
	"io"
	"io/ioutil"
)
import "gopkg.in/yaml.v2"

// YmlToJsonConverter converts from yml to json formats
type YmlToJsonConverter struct {
	prettyPrint bool
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
	tmpN := c.getNormalized(tmp)
	if c.prettyPrint {
		data, err = json.MarshalIndent(tmpN, "", "    ")
	} else {
		data, err = json.Marshal(tmpN)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (c *YmlToJsonConverter) WithPrettyPrint(prettyprint bool) *YmlToJsonConverter {
	c.prettyPrint = prettyprint
	return c
}

func (c *YmlToJsonConverter) getNormalized(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = c.getNormalized(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = c.getNormalized(v)
		}
	}
	return i
}
