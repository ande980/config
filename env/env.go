package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Provider is a type that implements config.Provider. The
// default env prefix is an empty string so the zero value
// is useful.
type Provider struct {
	prefix string
}

// New instantiates an empty usable Provider instance.
func New() *Provider {
	return &Provider{}
}

// WithPrefix allows for a custom application prefix
// for specified environmental variables.
func WithPrefix(prefix string) *Provider {
	return &Provider{prefix}
}

// Parse satisfies the config.Provider interface.
func (p *Provider) Parse(i interface{}) error {
	v := reflect.ValueOf(i)
	v = v.Elem()

	if !v.IsValid() {
		return nil
	}

	p.visit(v, p.prefix)

	return nil
}

func (p *Provider) visit(v reflect.Value, prefix string) {
	if v.Kind() != reflect.Struct {
		return
	}

	if !v.IsValid() {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		field = reflect.Indirect(field)

		name := v.Type().Field(i).Tag.Get("env")
		if name == "-" {
			continue
		}
		if name == "" {
			name = v.Type().Field(i).Name
			if prefix != "" {
				name = prefix + "_" + name
			}
		}
		name = strings.ToUpper(name)

		if field.Kind() == reflect.Struct {
			p.visit(field, name)
			continue
		}

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		usage := v.Type().Field(i).Tag.Get("usage")
		if usage == "" {
			// TODO
		}

		val := os.Getenv(name)
		if val == "" {
			continue
		}

		// Special case - has to go first or it clashes with *int64
		if field.Type() == reflect.TypeOf(time.Second) {
			dur, err := time.ParseDuration(val)
			if err != nil {
				panic(fmt.Errorf("parsing duration: %v", err))
			}
			field.Set(reflect.ValueOf(dur))
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			toCorrectIntType(field, val)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			toCorrectUintType(field, val)
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				panic(fmt.Errorf("parsing bool: %v", err))
			}
			field.SetBool(b)
		case reflect.Float64:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panic(fmt.Errorf("parsing float: %v", err))
			}
			field.SetFloat(f)
		}
	}
}

func toCorrectIntType(v reflect.Value, s string) {
	ival, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(fmt.Errorf("parsing int: %v", err))
	}

	switch v.Kind() {
	case reflect.Int:
		v.Set(reflect.ValueOf(int(ival)))
	case reflect.Int8:
		v.Set(reflect.ValueOf(int8(ival)))
	case reflect.Int16:
		v.Set(reflect.ValueOf(int16(ival)))
	case reflect.Int32:
		v.Set(reflect.ValueOf(int32(ival)))
	case reflect.Int64:
		v.SetInt(ival)
	}
	panic(&reflect.ValueError{Method: "toCorrectIntType", Kind: v.Kind()})
}

func toCorrectUintType(v reflect.Value, s string) {
	ival, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(fmt.Errorf("parsing int: %v", err))
	}

	switch v.Kind() {
	case reflect.Uint:
		v.Set(reflect.ValueOf(uint(ival)))
	case reflect.Uint8:
		v.Set(reflect.ValueOf(uint8(ival)))
	case reflect.Uint16:
		v.Set(reflect.ValueOf(uint16(ival)))
	case reflect.Uint32:
		v.Set(reflect.ValueOf(uint32(ival)))
	case reflect.Uint64:
		v.SetUint(ival)
	}
	panic(&reflect.ValueError{Method: "toCorrectUintType", Kind: v.Kind()})
}
