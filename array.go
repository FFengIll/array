package array

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotAStructPtr = errors.New("must spec a struct ptr to parse into")
var ErrUnsupportField = errors.New("can not parse unsupport field type (e.g. struct, ptr, ...)")
var ErrParse = errors.New("value parse error")
var ErrIndex = errors.New("array index invalid")

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

// doParse will parse the string array into ref (struct ptr) by tag like `array:"[1]"`.
//
// The array tag works like array element deref, e.g. `array:"[1]"` means array[1].
// Furthermore, add `omit` to allow to ignore absent items.
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
			return ErrUnsupportField
		default:
			tagContent, exist := tag.Lookup("array")
			// fmt.Println(config, exist)

			if exist {
				index, options := parseOptions(tagContent)
				if index < 0 || index >= len(data) {
					if options&OptOmitEmpty == OptOmitEmpty {
						continue
					} else {
						return ErrIndex
					}
				}

				// fmt.Println(data[index])
				parser, ok := defaultBuiltInParsers[typee.Kind()]
				if !ok {
					return ErrUnsupportField
				}
				value, err := parser(data[index])
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(value).Convert(typee))
			}
			break
		}
	}

	return nil
}

const (
	OptOmitEmpty int = 1 << iota
	OptUnknown
)

func parseOptions(value string) (index int, options int) {
	items := strings.Split(value, ",")
	if len(items) <= 0 {
		index = -1
		return
	}
	for _, it := range items {
		it := strings.TrimSpace(it)
		switch it {
		case "omitempty":
			options = options | OptOmitEmpty
		default:
			continue
		}
	}
	// first item must be like `[123]`
	indexStr := strings.TrimFunc(items[0], func(r rune) bool {
		if strings.ContainsRune("[] ", r) {
			return true
		}
		return false
	})
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		index = -1
	}
	return
}
