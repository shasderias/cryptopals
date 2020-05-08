package kv

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Marshal(v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil, errors.New("v must be a struct or a pointer to a struct")
	}

	buf := bytes.Buffer{}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		k := strings.ToLower(typ.Field(i).Name)
		var v string

		f := val.Field(i)
		switch f.Kind() {
		case reflect.String:
			v = f.String()
		case reflect.Int:
			v = strconv.Itoa(int(f.Int()))
		default:
			return nil, fmt.Errorf("error marshaling field %s, unsupported type %s", typ.Field(i).Name, f.Kind())
		}

		k = sanitizeKVString(k)
		v = sanitizeKVString(v)
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
		if i != typ.NumField()-1 {
			buf.WriteString("&")
		}
	}

	return buf.Bytes(), nil
}

func sanitizeKVString(s string) string {
	s = strings.ReplaceAll(s, "&", "_")
	s = strings.ReplaceAll(s, "=", "_")
	return s
}

func Unmarshal(data []byte, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return errors.New("v must be a non-nil pointer to a struct")
	}

	elem := val.Elem()
	elemType := val.Elem().Type()
	fields := map[string]int{}

	for i := 0; i < elemType.NumField(); i++ {
		fields[strings.ToLower(elemType.Field(i).Name)] = i
	}

	kvPairs := bytes.Split(data, []byte("&"))
	for i, kvPair := range kvPairs {
		s := bytes.Split(kvPair, []byte("="))
		if len(s) != 2 {
			return fmt.Errorf("key-value pair at index %d does not have exactly one = delimiter", i)
		}
		fieldIndex, ok := fields[strings.ToLower(string(s[0]))]
		if !ok {
			continue
		}

		switch elem.Field(fieldIndex).Kind() {
		case reflect.String:
			elem.Field(fieldIndex).SetString(string(s[1]))
		case reflect.Int:
			v, err := strconv.Atoi(string(s[1]))
			if err != nil {
				return fmt.Errorf("error unmarshaling value '%s' at index %d into int", string(s[1]), i)
			}
			elem.Field(fieldIndex).SetInt(int64(v))
		default:
			return fmt.Errorf("unsupported type '%s' at index %d", elem.Field(fieldIndex).Kind(), i)
		}
	}
	return nil
}
