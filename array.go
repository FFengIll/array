package array

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrNotAStructPtr = errors.New("")

type ParserFunc func(v string) (interface{}, error)

var defaultBuiltInParsers = map[reflect.Kind]ParserFunc{
	reflect.Bool: func(v string) (interface{}, error) {
		return strconv.ParseBool(v)
	},
	reflect.String: func(v string) (interface{}, error) {
		return v, nil
	},
	reflect.Int: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	},
	reflect.Int16: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	},
	reflect.Int32: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	},
	reflect.Int64: func(v string) (interface{}, error) {
		return strconv.ParseInt(v, 10, 64)
	},
	reflect.Int8: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	},
	reflect.Uint: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	},
	reflect.Uint16: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	},
	reflect.Uint32: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	},
	reflect.Uint64: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 64)
		return i, err
	},
	reflect.Uint8: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	},
	reflect.Float64: func(v string) (interface{}, error) {
		return strconv.ParseFloat(v, 64)
	},
	reflect.Float32: func(v string) (interface{}, error) {
		f, err := strconv.ParseFloat(v, 32)
		return float32(f), err
	},
}

func Parse(data []string, v interface{}) error {
	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotAStructPtr
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotAStructPtr
	}

	return doParse(data, ref)
}

func doParse(data []string, ref reflect.Value) error {
	refType := ref.Type()

	for i := 0; i < refType.NumField(); i++ {
		field := ref.Field(i)
		typeField := refType.Field(i)
		// fmt.Printf("%d %v %v\n", i, field, typeField)

		tag := typeField.Tag
		// fmt.Println(tag)

		typee := field.Type()
		switch typee.Kind() {
		case reflect.Ptr, reflect.Struct:
			return ErrNotAStructPtr
		default:
			value, exist := tag.Lookup("array")
			// fmt.Println(value, exist)

			if exist {
				value = value[1 : len(value)-1]
				index, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				if 0 <= index && index < len(data) {
					// fmt.Println(data[index])
					parser, ok := defaultBuiltInParsers[typee.Kind()]
					if !ok {
						return ErrNotAStructPtr
					}
					value, err := parser(data[index])
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(value).Convert(typee))
				}
			}
			break
		}
	}

	return nil
}
