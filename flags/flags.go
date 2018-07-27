package flags

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

// FlagSet embeds a flags.FlagSet for the purposes of defining
// flags at runtime using reflection and the struct passed
// to the Parse function.
type FlagSet struct {
	*flag.FlagSet
}

// New instantiates an empty usable flagset ready for parsing.
func New() *FlagSet {
	filepathNoExt := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	return &FlagSet{flag.NewFlagSet(filepathNoExt, flag.ContinueOnError)}
}

// Parse implements the config.Provider interface.
func (f *FlagSet) Parse(i interface{}) error {
	return f.parse(i, os.Args[1:]...)
}

func (f *FlagSet) parse(i interface{}, args ...string) error {
	v := reflect.ValueOf(i)
	v = v.Elem()

	if !v.IsValid() {
		return nil
	}

	f.visit(v, "")

	if err := f.FlagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return fmt.Errorf("parsing flags: %v", err)
	}

	return nil
}

func (f *FlagSet) visit(v reflect.Value, prefix string) {
	if v.Kind() != reflect.Struct {
		return
	}

	if !v.IsValid() {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		field = reflect.Indirect(field)

		name := v.Type().Field(i).Tag.Get("flag")
		if name == "-" {
			continue
		}
		if name == "" {
			name = v.Type().Field(i).Name
			if prefix != "" {
				name = prefix + "-" + name
			}
		}
		name = strings.ToLower(name)

		if field.Kind() == reflect.Struct {
			f.visit(field, name)
			continue
		}

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		usage := v.Type().Field(i).Tag.Get("usage")
		if usage == "" {
			// TODO
		}

		// Special case - has to go first or it clashes with *int64
		if field.Type() == reflect.TypeOf(time.Second) {
			f.DurationVar(field.Addr().Interface().(*time.Duration), name, time.Duration(field.Int()), usage)
			continue
		}

		switch field.Kind() {
		case reflect.String:
			f.StringVar(field.Addr().Interface().(*string), name, field.String(), usage)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.Int64Var(toInt64(field), name, field.Int(), usage)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			f.Uint64Var(toUint64(field), name, field.Uint(), usage)
		case reflect.Bool:
			f.BoolVar(field.Addr().Interface().(*bool), name, field.Bool(), usage)
		case reflect.Float64:
			f.Float64Var(field.Addr().Interface().(*float64), name, field.Float(), usage)
		}
	}
}

func toInt64(v reflect.Value) *int64 {
	switch v.Kind() {
	case reflect.Int:
		i := v.Addr().Interface().(*int)
		i64 := int64(*i)
		return &i64
	case reflect.Int8:
		i := v.Addr().Interface().(*int8)
		i64 := int64(*i)
		return &i64
	case reflect.Int16:
		i := v.Addr().Interface().(*int16)
		i64 := int64(*i)
		return &i64
	case reflect.Int32:
		i := v.Addr().Interface().(*int32)
		i64 := int64(*i)
		return &i64
	case reflect.Int64:
		return v.Addr().Interface().(*int64)
	}
	panic(&reflect.ValueError{Method: "toInt64", Kind: v.Kind()})
}

func toUint64(v reflect.Value) *uint64 {
	switch v.Kind() {
	case reflect.Uint:
		i := v.Addr().Interface().(*uint)
		ui64 := uint64(*i)
		return &ui64
	case reflect.Uint8:
		i := v.Addr().Interface().(*uint8)
		ui64 := uint64(*i)
		return &ui64
	case reflect.Uint16:
		i := v.Addr().Interface().(*uint16)
		ui64 := uint64(*i)
		return &ui64
	case reflect.Uint32:
		i := v.Addr().Interface().(*uint32)
		ui64 := uint64(*i)
		return &ui64
	case reflect.Uint64:
		return v.Addr().Interface().(*uint64)
	}
	panic(&reflect.ValueError{Method: "toUint64", Kind: v.Kind()})
}
