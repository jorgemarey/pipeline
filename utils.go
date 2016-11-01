package main

import (
	"bytes"
	"os"
	"strings"
	"text/template"
)

func parse(value interface{}, context map[string]interface{}) (interface{}, error) {
	switch typedValue := value.(type) {
	case string:
		b := &bytes.Buffer{}
		tmpl, err := template.New("parse-value").Parse(typedValue)
		if err = tmpl.Execute(b, context); err != nil {
			return typedValue, err
		}
		return b.String(), nil
	case map[string]interface{}:
		tmp := make(map[string]interface{})
		for key, innerValue := range typedValue {
			parsedValue, err := parse(innerValue, context)
			if err != nil {
				return tmp, err
			}
			tmp[key] = parsedValue
		}
		return tmp, nil
	case []interface{}:
		tmp := make([]interface{}, len(typedValue))
		for index, innerValue := range typedValue {
			parsedValue, err := parse(innerValue, context)
			if err != nil {
				return tmp, err
			}
			tmp[index] = parsedValue
		}
		return tmp, nil
	default:
		return typedValue, nil
	}
}

func GetEnvironment() map[string]string {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	return getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})
}
